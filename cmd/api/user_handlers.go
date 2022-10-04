package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/GarmaTs/linkshortener/internal/data"
	"github.com/GarmaTs/linkshortener/internal/validator"
	"github.com/google/uuid"
)

func (app *application) Welcome(w http.ResponseWriter, r *http.Request) {
	session, _, err := app.checkAuthorization(w, r)

	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf("Welcome %s!", session.Username)))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) LogOut(w http.ResponseWriter, r *http.Request) {
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
}

func (app *application) Signin(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"name"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &creds)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	u, err := app.models.Users.GetByName(creds.Username)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	match, err := u.Password.Matches(creds.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialResponse(w, r)
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)
	s := app.models.Sessions.Set(sessionToken, u.Name, expiresAt)
	if s.Username == "" {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:  input.Name,
		Email: input.Email,
	}
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateName):
			v.AddError("name", "a user with this name address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
