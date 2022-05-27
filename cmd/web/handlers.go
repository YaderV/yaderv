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
	var form userSignForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must a least 8 characters long")

	app.render(w, http.StatusOK, "home.tmpl")
}

func (app application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "home.tmpl")
}
