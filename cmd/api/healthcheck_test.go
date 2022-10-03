package main

import (
	"net/http"
	"testing"

	"github.com/GarmaTs/linkshortener/internal/assert"
)

func TestHealthcheckHandler(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, _ := ts.get(t, "/v1/healthcheck")

	assert.Equal(t, code, http.StatusOK)
}
