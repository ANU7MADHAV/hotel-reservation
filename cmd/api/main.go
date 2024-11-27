package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.23.3"

type config struct {
	port int
	env  string
	dsn  Dsn
}

type Dsn struct {
	dsn string
}

type Applications struct {
	config *config
	logger *log.Logger
}

func main() {
	var config config

	flag.IntVar(&config.port, "port", 4000, "API SERVER PORT")
	flag.StringVar(&config.env, "environment", "development", "ENVIRONMENT")
	flag.StringVar(&config.dsn.dsn, "db-dsn", "postgres://hotel::your_password@localhost/hotel_reservation?sslmode=disable", "POSTGRES SQL DATABASE DSN")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDb(&config)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	app := &Applications{
		config: &config,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.port),
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Printf("Starting %s server on %s ", config.env, srv.Addr)
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func openDb(cfg *config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.dsn.dsn)
	if err != nil {
		return nil, nil
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
