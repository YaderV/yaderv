package main

import (
	"context"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/YaderV/yaderv/internal/models"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
)

const templateRoot string = "./ui/html"
const staticPath string = "./ui/static"

// SessionUserIDKey represents the session where we are going to store
// the user id
const SessionUserIDKey string = "authenticatedUserID"

type application struct {
	infoLog        *log.Logger
	errorLog       *log.Logger
	templateCache  map[string]*template.Template
	users          models.UserModel
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// Get variables
	addr := flag.String("addr", ":4000", "Http network address")
	dns := flag.String("db-dns", os.Getenv("DB_DNS"), "DNS Database")
	flag.Parse()

	// Set custom loggers
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	db, err := openDB(*dns)

	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// We set a form decoder as helper to parse form values into
	// structs
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Tie some varibles to the application struct so we share data/fuctions
	// between the package istead of using global variables
	app := application{
		infoLog:        infoLog,
		errorLog:       errorLog,
		templateCache:  templateCache,
		users:          models.UserModel{DB: db},
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dns)

	if err != nil {
		return nil, err
	}

	// We create a context with a timeout of 5 seconds
	// to entablish a db pool connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verify/entablilsh the connection with the database
	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
