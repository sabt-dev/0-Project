package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIv      string
	Port      string
	DBHost    string
	JWTSecret string
	SMTPHost  string
	SMTPPort  string
	SMTPUser  string
	SMTPPass  string
	FromEmail string
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}

	AppConfig = Config{
		APIv:      os.Getenv("API_VERSION"),
		Port:      os.Getenv("PORT"),
		DBHost:    os.Getenv("DB_HOST"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		SMTPHost:  os.Getenv("SMTP_HOST"),
		SMTPPort:  os.Getenv("SMTP_PORT"),
		SMTPUser:  os.Getenv("SMTP_USER"),
		SMTPPass:  os.Getenv("SMTP_PASS"),
		FromEmail: os.Getenv("FROM_EMAIL"),
	}
}
