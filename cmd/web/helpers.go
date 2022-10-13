package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/GarmaTs/linkshortener/internal/data"
	"github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	_, err = buf.WriteTo(w)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) readShortUrlParam(r *http.Request) (string, error) {
	params := httprouter.ParamsFromContext(r.Context())
	shortUrl := params.ByName("short_url")
	if len(shortUrl) == 0 {
		return "", errors.New("invalid shortUrl parameter")
	}
	return shortUrl, nil
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}

func (app *application) checkAuthorization(w http.ResponseWriter, r *http.Request) (data.Session, string, error) {
	var session data.Session
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return session, "", data.ErrUnauthorized
		}
		return session, "", data.ErrUnauthorized
	}
	token := c.Value

	session, err = app.models.Sessions.Get(token)
	if err != nil {
		return session, "", data.ErrUnauthorized
	}
	return session, token, nil
}
