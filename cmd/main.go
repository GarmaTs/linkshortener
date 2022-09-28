package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	l5 "github.com/GarmaTs/linkshortener/internal/lesson5"
)

func main() {
	cfg := l5.Config{
		Version: "1.0.0",
		Port:    8000,
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := l5.Application{
		Config: cfg,
		Logger: logger,
	}
	app.FillFakeItems() // Fill items with fake data
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: app.Routes(),
	}

	logger.Printf("starting server on %s", srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
