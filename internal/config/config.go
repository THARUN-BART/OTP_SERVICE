package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all the environment and app settings
type Config struct {
	Port             string
	BrevoApiKey      string
	BrevoSenderEmail string
	OTPExpiryMinutes int
}

// LoadConfig reads the .env file and environment variables
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found or already loaded into system environment")
	}

	expiryMinutes, err := strconv.Atoi(getEnv("OTP_EXPIRY_MINUTES", "5"))
	if err != nil {
		expiryMinutes = 5
	}

	// Render sets the PORT environment variable automatically
	port := getEnv("PORT", getEnv("APP_PORT", "3000"))

	return &Config{
		Port:             port,
		BrevoApiKey:      getEnv("BREVO_API_KEY", ""),
		BrevoSenderEmail: getEnv("BREVO_SENDER_EMAIL", ""),
		OTPExpiryMinutes: expiryMinutes,
	}
}

// getEnv gets an environment variable or returns a fallback default
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
