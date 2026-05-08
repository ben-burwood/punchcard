package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	Port         int
	DBPath       string
	APIKey       string
	DashUser     string
	DashPassword string
	StaticDir    string
}

func Load() (Config, error) {
	cfg := Config{
		Port:         8080,
		DBPath:       getenv("DB_PATH", "/data/punchcard.db"),
		APIKey:       os.Getenv("PUNCHCARD_API_KEY"),
		DashUser:     os.Getenv("DASHBOARD_USER"),
		DashPassword: os.Getenv("DASHBOARD_PASSWORD"),
		StaticDir:    getenv("STATIC_DIR", "./frontend/dist"),
	}

	if v := os.Getenv("PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return cfg, errors.New("PORT must be an integer")
		}
		cfg.Port = p
	}

	if cfg.APIKey == "" {
		return cfg, errors.New("PUNCHCARD_API_KEY environment variable is required")
	}
	if cfg.DashUser == "" || cfg.DashPassword == "" {
		return cfg, errors.New("DASHBOARD_USER and DASHBOARD_PASSWORD environment variables are required")
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
