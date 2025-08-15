package config

import (
	"errors"
	"os"
	"time"
)

const (
	NotFoundHOST = "not found HOST in environment variables"
	NotFoundUSER = "not found USER in environment variables"
	NotFoundPASS = "not found PASSWORD in environment variables"
	NotFoundDB   = "not found DB_NAME in environment variables"
)

type Config struct {
	ListenPort string

	Port     string
	Host     string
	User     string
	Password string
	DB       string
	Interval time.Duration
}

func NewConfig() (*Config, error) {
	return readConfig()
}

func readConfig() (*Config, error) {
	listenPort := os.Getenv("LISTEN_PORT")
	if listenPort == "" {
		listenPort = "5454"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "5432"
	}
	host := os.Getenv("HOST")
	if host == "" {
		return nil, errors.New(NotFoundHOST)
	}
	user := os.Getenv("USER")
	if user == "" {
		return nil, errors.New(NotFoundUSER)
	}
	password := os.Getenv("PASSWORD")
	if password == "" {
		return nil, errors.New(NotFoundPASS)
	}
	db := os.Getenv("DB_NAME")
	if db == "" {
		return nil, errors.New(NotFoundDB)
	}
	interval := os.Getenv("INTERVAL")
	if interval == "" {
		interval = "5s"
	}
	intervalDuration, err := time.ParseDuration(interval)
	if err != nil {
		return nil, err
	}

	return &Config{
		ListenPort: listenPort,
		Port:       port,
		Host:       host,
		User:       user,
		Password:   password,
		DB:         db,
		Interval:   intervalDuration,
	}, nil
}
