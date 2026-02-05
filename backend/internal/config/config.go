package config

import (
	"os"
)

const (
	envGoogleClientID     = "GOOGLE_CLIENT_ID"
	envGoogleClientSecret = "GOOGLE_CLIENT_SECRET"
	envSessionSecret      = "SESSION_SECRET"
	envBackendURL         = "BACKEND_URL"
	envFrontendURL        = "FRONTEND_URL"
)

type Config struct {
	GoogleClientID     string
	GoogleClientSecret string
	SessionSecret      string
	BackendURL         string
	FrontendURL        string
}

func Load() Config {
	return Config{
		GoogleClientID:     os.Getenv(envGoogleClientID),
		GoogleClientSecret: os.Getenv(envGoogleClientSecret),
		SessionSecret:      getEnvOrDefault(envSessionSecret, "dev-session-secret-change-in-production"),
		BackendURL:         getEnvOrDefault(envBackendURL, "http://localhost:8080"),
		FrontendURL:        getEnvOrDefault(envFrontendURL, "http://localhost:3000"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
