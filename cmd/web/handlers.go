package main

import (
	"net/http"

	"github.com/YaderV/yaderv/internal/validator"
)

type userSignForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app application) userSignup(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "signup.tmpl")
}

func (app application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	app.render(w, http.StatusOK, "home.tmpl")
}

func (app application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "home.tmpl")
}
