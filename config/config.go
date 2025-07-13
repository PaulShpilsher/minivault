package config

import (
	"os"
)

type Config struct {
	ServerPort  string
	OllamaURL   string
	OllamaModel string
}

func Load() *Config {
	cfg := &Config{
		ServerPort:  getEnv("MINIVAULT_PORT", ":8080"),
		OllamaURL:   getEnv("OLLAMA_URL", "http://localhost:11434/api/chat"),
		OllamaModel: getEnv("OLLAMA_MODEL", "gemma:2b"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
