package main

import (
	"net/http"

	"github.com/GarmaTs/linkshortener/internal/data"
	"github.com/GarmaTs/linkshortener/internal/validator"
)

func (app *application) addUrlByUserNameAndFullUrl(w http.ResponseWriter, r *http.Request) {
	session, _, err := app.checkAuthorization(w, r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}
	var input struct {
		FullUrl string `json:"full_url"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	url := &data.Url{
		FullUrl: input.FullUrl,
	}

	v := validator.New()
	if data.ValidateUrl(v, url); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Urls.Insert(url, session.Username, input.FullUrl)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"url": url}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getFullUrlByShortUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl, err := app.readShortUrlParam(r)
	if err != nil || len(shortUrl) == 0 {
		app.notFoundResponse(w, r)
		return
	}
	url := data.Url{
		ShortUrl: shortUrl,
	}

	err = app.models.Urls.GetOne(&url, shortUrl)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"url": url.FullUrl}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
