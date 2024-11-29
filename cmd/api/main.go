package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"hotel-reservation/internal/data"
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
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type Applications struct {
	config *config
	logger *log.Logger
	data   data.Models
}

func main() {
	var config config

	flag.IntVar(&config.port, "port", 4000, "API SERVER PORT")
	flag.StringVar(&config.env, "environment", "development", "ENVIRONMENT")
	flag.StringVar(&config.dsn.dsn, "db-dsn", "postgres://hotel:your_password@localhost:5432/hotel_reservation?sslmode=disable", "POSTGRES SQL DATABASE DSN")
	flag.IntVar(&config.dsn.maxOpenConns, "db-max-open-conns", 25, "POSTGRES maximum open connections")
	flag.IntVar(&config.dsn.maxIdleConns, "db-max-idle-conns", 25, "POSTGRES maximum idle connections")
	flag.StringVar(&config.dsn.maxIdleTime, "db-max-idle-time", "15m", "POSTGRES maximum idle time")

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
		data:   data.NewHotelModel(db),
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
		return nil, err
	}
	log.Println("Setting connection pool parameters...")

	db.SetMaxOpenConns(cfg.dsn.maxOpenConns)
	db.SetMaxIdleConns(cfg.dsn.maxIdleConns)

	duration, err := time.ParseDuration(cfg.dsn.maxIdleTime)

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
	fmt.Println("db", db)

	return db, nil
}
