package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/GarmaTs/linkshortener/internal/assert"
	"github.com/GarmaTs/linkshortener/internal/data"
)

func TestWelcome(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, _ := ts.get(t, "/v1/welcome")
	assert.Equal(t, code, http.StatusUnauthorized)
}

func TestSigninStatusOk(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := &application{
		logger: logger,
		models: data.FakeNewModels(),
	}

	input := struct {
		Name     string
		Email    string
		Password string
	}{
		"testname",
		"test@test.test",
		"testpa55word",
	}
	user := &data.User{
		Name:  input.Name,
		Email: input.Email,
	}
	user.ID = 1
	user.CreatedAt = time.Now()
	err := user.Password.Set(input.Password)
	if err != nil {
		t.Fatal(err)
	}
	err = app.models.Users.Insert(user)
	if err != nil {
		t.Fatal(err)
	}

	creds := struct {
		Username string `json:"name"`
		Password string `json:"password"`
	}{
		input.Name,
		input.Password,
	}

	b := &bytes.Buffer{}
	err = json.NewEncoder(b).Encode(creds)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "/", b)
	if err != nil {
		t.Fatal(err)
	}

	app.Signin(rr, r)
	rs := rr.Result()
	defer rs.Body.Close()
	assert.Equal(t, rs.StatusCode, http.StatusOK)
}

func TestSigninStatusAnauthorized(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := &application{
		logger: logger,
		models: data.FakeNewModels(),
	}

	input := struct {
		Name     string
		Email    string
		Password string
	}{
		"testname",
		"test@test.test",
		"testpa55word",
	}
	user := &data.User{
		Name:  input.Name,
		Email: input.Email,
	}
	user.ID = 1
	user.CreatedAt = time.Now()
	err := user.Password.Set(input.Password)
	if err != nil {
		t.Fatal(err)
	}
	err = app.models.Users.Insert(user)
	if err != nil {
		t.Fatal(err)
	}

	creds := struct {
		Username string `json:"name"`
		Password string `json:"password"`
	}{
		input.Name,
		"wrongpassword",
	}

	b := &bytes.Buffer{}
	err = json.NewEncoder(b).Encode(creds)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "/", b)
	if err != nil {
		t.Fatal(err)
	}

	app.Signin(rr, r)
	rs := rr.Result()
	defer rs.Body.Close()
	assert.Equal(t, rs.StatusCode, http.StatusUnauthorized)
}

func TestRegisterUserHandler(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := &application{
		logger: logger,
		models: data.FakeNewModels(),
	}

	creds := struct {
		Username string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		"newtestuser",
		"newtestuser@test.com",
		"testuserpa55word",
	}

	b := &bytes.Buffer{}
	err := json.NewEncoder(b).Encode(creds)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "/", b)
	defer r.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	app.registerUserHandler(rr, r)

	_, err = app.models.Users.GetByName(creds.Username)
	if err != nil {
		t.Errorf("User was not added")
	}
}
