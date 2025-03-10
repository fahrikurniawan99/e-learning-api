package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbHost             string
	DbUser             string
	DbPassword         string
	DbName             string
	DbPort             string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectUri  string
	WaWebUrl           string
	AppEnv             string
}

var Env *Config

func LoadEnv() {
	// Load .env jika ada
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	Env = &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		DbHost:             getEnv("DB_HOST", "localhost"),
		DbUser:             getEnv("DB_USER", "root"),
		DbPassword:         getEnv("DB_PASSWORD", ""),
		DbName:             getEnv("DB_NAME", "db_name"),
		DbPort:             getEnv("DB_PORT", "5432"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectUri:  getEnv("GOOGLE_REDIRECT_URI", ""),
		WaWebUrl:           getEnv("WA_WEB_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
