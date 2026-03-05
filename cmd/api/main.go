package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type app struct {
	logger *log.Logger
	config config
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "port to host the api")
	flag.StringVar(&cfg.env, "env", "development", "environment of the code could be (production|development|staging)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &app{
		logger: logger,
		config: cfg,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("server started on port %d in %s environment", cfg.port, cfg.env)
	err := server.ListenAndServe()
	logger.Fatal(err)
}
