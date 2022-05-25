package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app application) routes() http.Handler {
	router := httprouter.New()
	fileServer := http.FileServer(http.Dir(staticPath))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Users
	router.HandlerFunc(http.MethodGet, "/user/signup/", app.userSignup)
	router.HandlerFunc(http.MethodPost, "/user/signup/", app.userSignupPost)

	router.HandlerFunc(http.MethodGet, "/", app.home)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
