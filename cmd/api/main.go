package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GarmaTs/linkshortener/internal/data"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

type myData struct {
	Conf struct {
		Port int
		Env  string
		Db   struct {
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
		port: md.Conf.Port,
		env:  md.Conf.Env,
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

	app := &application{
		config: *cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	defer srv.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM)

	go func() {
		logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
		err := srv.ListenAndServe()
		logger.Fatal(err)
	}()

	<-stop
	fmt.Println("Server stopped...")
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
