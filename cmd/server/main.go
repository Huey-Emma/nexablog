package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"nexablog/config"
	"nexablog/db"
	"nexablog/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.New(db.Config{
		Uri:       cfg.DB.Uri,
		IdleConns: cfg.DB.IdleConns,
		OpenConns: cfg.DB.OpenConns,
		IdleTime:  cfg.DB.IdleTime,
	})
	if err != nil {
		log.Fatal(err)
	}

	appl := app.New(cfg, database)

	signals := []os.Signal{os.Interrupt, os.Kill}

	ctx, cancel := signal.NotifyContext(context.Background(), signals...)
	defer cancel()

	if err := appl.StartAndRun(ctx); err != nil {
		log.Fatal(err)
	}
}
