package config

import (
	"os"
)

type Config struct {
	WSURL string
}

func Load() Config {
	url := os.Getenv("SKYBIT_WS_URL")
	if url == "" {
		// Default to local dev; replace as needed
		url = "ws://localhost:8080/ws"
	}
	return Config{WSURL: url}
}
