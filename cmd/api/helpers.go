package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/GarmaTs/linkshortener/internal/data"
	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}

func (app *application) readShortUrlParam(r *http.Request) (string, error) {
	params := httprouter.ParamsFromContext(r.Context())

	shortUrl := params.ByName("short_url")
	if len(shortUrl) == 0 {
		return "", errors.New("invalid shortUrl parameter")
	}
	return shortUrl, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
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
