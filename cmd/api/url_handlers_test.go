package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/GarmaTs/linkshortener/internal/assert"
	"github.com/GarmaTs/linkshortener/internal/data"
)

func TestGetFullUrlByShortUrl(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	url := data.Url{
		ID:        1,
		CreatedAt: time.Now(),
		ShortUrl:  "short",
		FullUrl:   "full",
	}
	err := app.models.Urls.Insert(&url, "", url.FullUrl)
	if err != nil {
		t.Fatal(err)
	}

	code, _, _ := ts.get(t, "/v1/urls/"+url.ShortUrl)

	assert.Equal(t, code, http.StatusOK)
}
