package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"time"

	"github.com/Dagime-Teshome/greenlight/internal/data"
	"github.com/Dagime-Teshome/greenlight/internal/jsonlog"
	_ "github.com/lib/pq"
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

	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type app struct {
	logger *jsonlog.Logger
	config config
	models data.Models
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "port to host the api")
	flag.StringVar(&cfg.env, "env", "development", "environment of the code could be (production|development|staging)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rqs-number", 2, "Request Per Second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst-limit", 4, "Limit for Request Burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Flag for enabling/disabling rate limiter")
	flag.Parse()

	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)
	app := &app{
		logger: logger,
		config: cfg,
		models: data.NewModel(db),
	}
	err = serve(app)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {

	db, err := sql.Open("postgres", cfg.db.dsn)
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
	return db, nil
}
