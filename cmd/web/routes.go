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

	// Set dynamic chain only for handlers that require sessions
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// Users
	router.Handler(http.MethodGet, "/user/signup/", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup/", dynamic.ThenFunc(app.userSignupPost))

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
