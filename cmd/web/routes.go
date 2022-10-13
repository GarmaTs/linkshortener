package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/url/:short_url", app.redirectToFullUrl)
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/signin", app.signinForm)
	router.HandlerFunc(http.MethodPost, "/signin", app.signinPost)
	router.HandlerFunc(http.MethodGet, "/signup", app.signupForm)
	router.HandlerFunc(http.MethodPost, "/signup", app.signupPost)
	router.HandlerFunc(http.MethodGet, "/urls", app.urlsForm)
	router.HandlerFunc(http.MethodPost, "/urls", app.urlsPost)
	router.HandlerFunc(http.MethodGet, "/links", app.linksForm)
	router.HandlerFunc(http.MethodPost, "/logout", app.logout)

	return app.recoverPanic(app.logRequest(router))
}
