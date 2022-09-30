package lesson5

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) Routes() *httprouter.Router {
	router := httprouter.New()
	router.RedirectTrailingSlash = false

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/items", app.showItemsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/items/:id", app.showSingleItemHandler)

	return router
}
