package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/GarmaTs/linkshortener/internal/data"
	"github.com/GarmaTs/linkshortener/internal/validator"
	"github.com/google/uuid"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}

func (app *application) signinForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSingInForm{}
	app.render(w, http.StatusOK, "signin.html", data)
}

type userSingInForm struct {
	Name                string `form:"name"`
	validator.Validator `form:"-"`
}

func (app *application) signinPost(w http.ResponseWriter, r *http.Request) {
	var form userSingInForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	password := r.PostForm.Get("password")
	form.Check(form.Name != "", "name", "Must be provided")
	form.Check(password != "", "password", "Must be provided")
	if !form.Valid() {
		app.showCheckedPage(w, r, form, "signin.html")
		return
	}

	u, err := app.models.Users.GetByName(form.Name)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			form.AddError("autherr", "Incorrect name or password")
			app.showCheckedPage(w, r, form, "signin.html")
			return
		} else {
			app.clientError(w, http.StatusInternalServerError)
			return
		}
	}

	match, err := u.Password.Matches(password)
	if err != nil {
		form.AddError("autherr", "Incorrect name or password")
		app.showCheckedPage(w, r, form, "signin.html")
		return
	}
	if !match || u == nil {
		form.AddError("autherr", "Incorrect name or password")
		app.showCheckedPage(w, r, form, "signin.html")
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(24 * time.Hour)
	s := app.models.Sessions.Set(sessionToken, u.Name, expiresAt)
	if s.Username == "" {
		app.serverError(w, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	http.Redirect(w, r, "/urls", http.StatusSeeOther)
}

func (app *application) showCheckedPage(w http.ResponseWriter, r *http.Request, form any, pageName string) {
	data := app.newTemplateData(r)
	data.Form = form
	app.render(w, http.StatusUnprocessableEntity, pageName, data)
}

func (app *application) redirectToFullUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl, err := app.readShortUrlParam(r)
	if err != nil || len(shortUrl) == 0 {
		app.notFound(w)
		return
	}
	url := data.Url{
		ShortUrl: shortUrl,
	}

	err = app.models.Urls.GetOne(&url, shortUrl)
	if err != nil {
		app.notFound(w)
		return
	}

	http.Redirect(w, r, url.FullUrl, http.StatusSeeOther)
}

type userSingUpForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password1           string `form:"password1"`
	Password2           string `form:"password2"`
	validator.Validator `form:"-"`
}

func (app *application) signupForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSingUpForm{}
	app.render(w, http.StatusOK, "signup.html", data)
}

func validateSignUp(form *userSingUpForm) {
	form.Check(form.Name != "", "name", "Must be provided")
	form.Check(form.Email != "", "email", "Must be provided")
	form.Check(form.Password1 != "", "password1", "Must be provided")
	form.Check(len(form.Name) >= 3, "name", "Min length for name is 3")
	form.Check(len(form.Password1) >= 5, "password1", "Min length for password is 5")
	form.Check(len(form.Name) <= 30, "name", "Max length for name is 30")
	form.Check(validator.Matches(form.Email, validator.EmailRX), "email", "Must be a valid email address")
	form.Check(len(form.Password1) <= 20, "Password1", "Max length for password is 20")
	form.Check(form.Password1 == form.Password2, "password1", "Provided passwords were different")
}

func (app *application) signupPost(w http.ResponseWriter, r *http.Request) {
	var form userSingUpForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	validateSignUp(&form)
	if !form.Valid() {
		app.showCheckedPage(w, r, form, "signup.html")
		return
	}

	user := &data.User{
		Name:  form.Name,
		Email: form.Email,
	}
	err = user.Password.Set(form.Password1)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateName):
			form.AddError("commonerr", "User with this name already exists, please choose another.")
			app.showCheckedPage(w, r, form, "signup.html")
		default:
			app.serverError(w, err)
		}
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(24 * time.Hour)
	s := app.models.Sessions.Set(sessionToken, user.Name, expiresAt)
	if s.Username == "" {
		app.serverError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	http.Redirect(w, r, "/urls", http.StatusSeeOther)
}

type urlsForm struct {
	FullUrl             string `form:"full_url"`
	validator.Validator `form:"-"`
	Result              map[string]string
}

func (app *application) urlsForm(w http.ResponseWriter, r *http.Request) {
	_, _, err := app.checkAuthorization(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	data := app.newTemplateData(r)
	data.Form = urlsForm{}
	app.render(w, http.StatusOK, "url_create.html", data)
}

func (app *application) urlsPost(w http.ResponseWriter, r *http.Request) {
	session, _, err := app.checkAuthorization(w, r)
	if err != nil {
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	var form urlsForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.Check(form.FullUrl != "", "full_url", "must be provided")
	form.Check(len(form.FullUrl) <= 1000, "full_url", "must be less than 1000 characters long")
	form.Check(strings.HasPrefix(form.FullUrl, "http"), "full_url", "has wrong format")

	if !form.Valid() {
		app.showCheckedPage(w, r, form, "url_create.html")
		return
	}

	url := &data.Url{
		FullUrl: form.FullUrl,
	}
	err = app.models.Urls.Insert(url, session.Username, form.FullUrl)
	if err != nil {
		app.clientError(w, http.StatusInternalServerError)
		return
	}

	form.Result = map[string]string{}

	form.Result["short_url"] = fmt.Sprintf("%s:%d%s%s", app.config.host, app.config.port, app.config.shortUrlPrefix, url.ShortUrl)
	data := app.newTemplateData(r)
	data.Form = form
	app.render(w, http.StatusOK, "url_create.html", data)
}

func (app *application) linksForm(w http.ResponseWriter, r *http.Request) {
	session, _, err := app.checkAuthorization(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	urls, err := app.models.Urls.GetList(session.Username)
	if err != nil {
		app.clientError(w, http.StatusInternalServerError)
		return
	}
	for i := 0; i < len(urls); i++ {
		// Since Url.ShortUrl is only last part of url, we need to concatenate with domain, port and url prefix
		urls[i].ShortUrl = fmt.Sprintf("%s:%d%s%s", app.config.host, app.config.port, app.config.shortUrlPrefix, urls[i].ShortUrl)
	}

	data := app.newTemplateData(r)
	data.Urls = urls
	data.Form = urlsForm{}
	app.render(w, http.StatusOK, "user_links.html", data)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	_, token, err := app.checkAuthorization(w, r)
	if err != nil {
		return
	}
	app.models.Sessions.Remove(token)
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}
