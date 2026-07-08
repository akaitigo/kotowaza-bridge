package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds application configuration.
type Config struct {
	Port               string
	DatabaseURL        string
	LLMAPIKey          string
	LLMModel           string
	CORSAllowedOrigins []string
}

// defaultCORSOrigin is the local development frontend origin used only when
// CORS_ALLOWED_ORIGINS is not set. Production deployments must set the
// environment variable explicitly to their real origins.
const defaultCORSOrigin = "http://localhost:3000"

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	llmKey := os.Getenv("LLM_API_KEY")
	if llmKey == "" {
		return nil, fmt.Errorf("LLM_API_KEY is required")
	}

	llmModel := os.Getenv("LLM_MODEL")
	if llmModel == "" {
		llmModel = "claude-sonnet-4-20250514"
	}

	return &Config{
		Port:               port,
		DatabaseURL:        dbURL,
		LLMAPIKey:          llmKey,
		LLMModel:           llmModel,
		CORSAllowedOrigins: parseCORSOrigins(os.Getenv("CORS_ALLOWED_ORIGINS")),
	}, nil
}

// parseCORSOrigins splits a comma-separated origin list, trimming whitespace
// and dropping empty entries. When no origin is provided it falls back to the
// local development origin so the app remains usable without configuration.
func parseCORSOrigins(raw string) []string {
	var origins []string
	for _, part := range strings.Split(raw, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	if len(origins) == 0 {
		return []string{defaultCORSOrigin}
	}
	return origins
}
