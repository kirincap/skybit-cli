package snaptrade

import (
    "fmt"
    "os"
)

type Config struct {
    ClientID     string
    ClientSecret string
    Env          string // sandbox|production
}

func LoadConfig() (Config, error) {
    id := os.Getenv("SNAPTRADE_CLIENT_ID")
    sec := os.Getenv("SNAPTRADE_CLIENT_SECRET")
    env := os.Getenv("SNAPTRADE_ENV")
    if env == "" { env = "sandbox" }
    if id == "" || sec == "" {
        return Config{}, fmt.Errorf("missing SNAPTRADE_CLIENT_ID/SECRET envs")
    }
    return Config{ClientID: id, ClientSecret: sec, Env: env}, nil
}


