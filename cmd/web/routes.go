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
	router.Handler(http.MethodGet, "/user/login/", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login/", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodGet, "/user/signup/", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup/", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodPost, "/user/logout/", dynamic.ThenFunc(app.userLogoutPost))

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))

	// private urls
	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/manage/article/", protected.ThenFunc(app.articleManageList))
	router.Handler(http.MethodGet, "/manage/article/create/", protected.ThenFunc(app.articleCreate))
	router.Handler(http.MethodPost, "/manage/article/create/", protected.ThenFunc(app.articleCreatePost))
	router.Handler(http.MethodGet, "/manage/article/edit/:id/", protected.ThenFunc(app.articleEdit))
	//router.Handler(http.MethodPost, "/manage/article/:id/edit/", protected.ThenFunc(app.articleCreate))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
