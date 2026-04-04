package config

import (
	"fmt"
	"os"
)

// Config holds application configuration.
type Config struct {
	Port        string
	DatabaseURL string
	LLMAPIKey   string
	LLMModel    string
	CORSOrigin  string
}

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

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "http://localhost:3000"
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		LLMAPIKey:   llmKey,
		LLMModel:    llmModel,
		CORSOrigin:  corsOrigin,
	}, nil
}
