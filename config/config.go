package config

import (
	"os"
	"strconv"
)

type Config struct {
	MongoURI           string
	DatabaseName       string
	SessionSecret      string
	SMTPHost           string
	SMTPPort           int
	SMTPUsername       string
	SMTPPassword       string
	FromEmail          string
	BaseURL            string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleCallbackURL  string
}

// NewConfig todo: Create .env file for these
func NewConfig() *Config {
	return &Config{
		MongoURI:           os.Getenv("MONGO_URI"),
		SMTPPort:           func() int { port, _ := strconv.Atoi(os.Getenv("SMTP_PORT")); return port }(),
		SessionSecret:      os.Getenv("SESSION_SECRET"),
		SMTPHost:           os.Getenv("SMTP_HOST"),
		SMTPUsername:       os.Getenv("SMTP_USERNAME"),
		SMTPPassword:       os.Getenv("SMTP_PASSWORD"),
		FromEmail:          os.Getenv("FROM_EMAIL"),
		BaseURL:            os.Getenv("BASE_URL"),
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleCallbackURL:  os.Getenv("GOOGLE_CALLBACK_URL"),
	}
}
