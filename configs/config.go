package main

import (
	"log"
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	
	// OAuth Configuration
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	
	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURL  string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "sso_app.db"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
		
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
		
		GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GitHubRedirectURL:  getEnv("GITHUB_REDIRECT_URL", "http://localhost:8080/auth/github/callback"),
	}
	
	// Validate required OAuth settings
	if config.GoogleClientID == "" {
		log.Println("Warning: GOOGLE_CLIENT_ID not set. Google OAuth will not work.")
	}
	if config.GitHubClientID == "" {
		log.Println("Warning: GITHUB_CLIENT_ID not set. GitHub OAuth will not work.")
	}
	
	return config
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
