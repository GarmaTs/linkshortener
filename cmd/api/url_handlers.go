package main

import (
	"net/http"

	"github.com/GarmaTs/linkshortener/internal/data"
)

func (app *application) getFullUrlByShortUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl, err := app.readShortUrlParam(r)
	if err != nil || len(shortUrl) == 0 {
		app.notFoundResponse(w, r)
		return
	}
	url := data.Url{
		ShortUrl: "test",
		FullUrl:  "testFullUrl",
	}
	if shortUrl != url.ShortUrl {
		app.notFoundResponse(w, r)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"url": url}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
