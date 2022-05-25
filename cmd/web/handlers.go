package main

import "net/http"

func (app application) userSignup(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "home.tmpl")
}

func (app application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "home.tmpl")
}

func (app application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "home.tmpl")
}
