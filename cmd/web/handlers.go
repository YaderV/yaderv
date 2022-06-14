package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/YaderV/yaderv/internal/models"
	"github.com/YaderV/yaderv/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type userSignForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	RepeatedPassword    string `form:"repeated_password"`
	validator.Validator `form:"-"`
}

func (app application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
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
	form.CheckField(validator.ConfirmPassword(form.Password, form.RepeatedPassword), "password", "Both password must the same")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "User Created")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl", data)
}

func (app application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)

	if err != nil {
		// With have a problem with the credentials, so we add a non field error
		// and renders the login page again
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or Password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
			return
		}
		app.serverError(w, err)
		return
	}

	// Generate a new session id
	err = app.sessionManager.RenewToken(r.Context())

	if err != nil {
		app.serverError(w, err)
		return
	}

	// We store a new user id in the session
	app.sessionManager.Put(r.Context(), SessionUserIDKey, id)
	app.sessionManager.Put(r.Context(), "flash", "The user has logged in succefully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), SessionUserIDKey)
	app.sessionManager.Put(r.Context(), "flash", "The user has been logged out succefully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app application) articleManageList(w http.ResponseWriter, r *http.Request) {
	articles, err := app.articles.List()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Articles = articles
	app.render(w, http.StatusOK, "article_list.tmpl", data)
}

type articleCreateForm struct {
	Title               string   `form:"title"`
	Body                string   `form:"body"`
	Categories          []string `form:"categories[]"`
	validator.Validator `form:"-"`
}

func (app application) articleCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = &articleCreateForm{}
	app.render(w, http.StatusOK, "article_create.tmpl", data)
}

func (app application) articleCreatePost(w http.ResponseWriter, r *http.Request) {
	var form articleCreateForm
	err := app.decodePostForm(r, &form)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Body), "body", "This field cannot be blank")
	form.CheckField(validator.NotEmpty(form.Categories), "categories", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusOK, "article_create.tmpl", data)
		return
	}

	err = app.articles.Create(form.Title, form.Body, form.Categories)

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Article Created")
	http.Redirect(w, r, "/manage/article", http.StatusSeeOther)
}

func (app application) articleEdit(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	article, err := app.articles.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Form = &articleCreateForm{
		Title:      article.Title,
		Body:       article.Body,
		Categories: article.Categories,
	}
	app.render(w, http.StatusOK, "article_create.tmpl", data)
}
