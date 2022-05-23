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
)

const templateRoot string = "./ui/html"

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
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

	if db != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Tie some varibles to the application struct so we share data/fuctions
	// between the package istead of using global variables
	app := application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: templateCache,
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
