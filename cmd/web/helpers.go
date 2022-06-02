package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/form/v4"
)

func (app application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app application) render(w http.ResponseWriter, status int, name string, data *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		err := fmt.Errorf("The template %s does not exist", name)
		app.serverError(w, err)
		return
	}

	// Exec the template and store it in a buffer to check if it's corrent
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Modify the header of the resposer
	w.WriteHeader(status)
	buf.WriteTo(w)
}

// decodePostForm wraps the logic of parsing a post request
// and mapping the values to a destination
func (app application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// the form Decoder requires a non-nil struct as destiniation
		// if this error happens we want to manage as a programming (app level)
		// error rather than a client (bad request) error, so we raise a panic
		// that shoud be treated like a server error
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

	}
	return nil
}

func (app application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), SessionUserIDKey)
}
