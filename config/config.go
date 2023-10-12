package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var ErrNoValue = errors.New("no value")

type Config struct {
	Port string
	DB   struct {
		Uri string
		IdleConns,
		OpenConns,
		IdleTime int
	}
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprint(8000)
	}

	uri := os.Getenv("PG_URI")
	if uri == "" {
		return nil, fmt.Errorf("PG_URI value is an empty string: %w", ErrNoValue)
	}

	idleconns, err := strconv.Atoi(os.Getenv("PG_IDLE_CONNS"))
	if err != nil {
		return nil, err
	}

	openconns, err := strconv.Atoi(os.Getenv("PG_OPEN_CONNS"))
	if err != nil {
		return nil, err
	}

	idletime, err := strconv.Atoi(os.Getenv("PG_IDLE_TIME"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Port: port,
		DB: struct {
			Uri string
			IdleConns,
			OpenConns,
			IdleTime int
		}{
			Uri:       uri,
			IdleConns: idleconns,
			OpenConns: openconns,
			IdleTime:  idletime,
		},
	}

	return cfg, nil
}
