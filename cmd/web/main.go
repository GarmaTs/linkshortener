package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/GarmaTs/linkshortener/internal/data"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

type config struct {
	port           int
	env            string
	host           string
	shortUrlPrefix string
	db             struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config        config
	logger        *log.Logger
	models        data.Models
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

type myData struct {
	Conf struct {
		Port           int
		Env            string
		Host           string
		ShortUrlPrefix string `yaml:"short_url_prefix"`
		Db             struct {
			Dsn          string `yaml:"dsn"`
			MaxOpenConns int    `yaml:"maxOpenConns"`
			MaxIdleConns int    `yaml:"maxIdleConns"`
			MaxIdleTime  string `yaml:"maxIdleTime"`
		}
	}
}

func readConf(filename string) (*config, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	md := &myData{}
	err = yaml.Unmarshal(buf, md)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	c := &config{
		port:           md.Conf.Port,
		env:            md.Conf.Env,
		host:           md.Conf.Host,
		shortUrlPrefix: md.Conf.ShortUrlPrefix,
		db: struct {
			dsn          string
			maxOpenConns int
			maxIdleConns int
			maxIdleTime  string
		}{
			md.Conf.Db.Dsn,
			md.Conf.Db.MaxOpenConns,
			md.Conf.Db.MaxIdleConns,
			md.Conf.Db.MaxIdleTime},
	}

	return c, err
}

func main() {
	cfg, err := readConf("./config/conf.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Printf("database connection pool established")

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	app := &application{
		config:        *cfg,
		logger:        logger,
		models:        data.NewModels(db),
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Printf("Starting server on %d", app.config.port)
	err = srv.ListenAndServe()
	app.logger.Fatal(err)
}

func openDB(cfg *config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, err
}
