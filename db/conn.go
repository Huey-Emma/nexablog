package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

type Config struct {
	Uri string
	IdleConns,
	OpenConns,
	IdleTime int
}

func New(cfg Config) (*DB, error) {
	dsn, err := pq.ParseURL(cfg.Uri)
	if err != nil {
		return nil, err
	}

	client, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	client.SetMaxIdleConns(cfg.IdleConns)
	client.SetMaxOpenConns(cfg.OpenConns)
	client.SetConnMaxIdleTime(time.Duration(cfg.IdleTime) * time.Second)

	db := &DB{client}

	return db, nil
}

func (db *DB) PingDB(ctx context.Context) error {
	return db.PingContext(ctx)
}

func (db *DB) CloseDB() error {
	return db.Close()
}
