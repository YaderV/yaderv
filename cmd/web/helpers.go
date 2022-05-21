package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
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

func (app application) render(w http.ResponseWriter, status int, name string) {
	ts, ok := app.templateCache[name]
	if !ok {
		err := fmt.Errorf("The template %s does not exist", name)
		app.serverError(w, err)
		return
	}

	// Exec the template and store it in a buffer to check if it's corrent
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", nil)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Modify the header of the resposer
	w.WriteHeader(status)
	buf.WriteTo(w)
}
