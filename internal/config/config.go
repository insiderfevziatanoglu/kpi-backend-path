package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	AppEnv     string
	DBUrl      string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Env not found")
	}

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		AppEnv:     getEnv("APP_ENV", "development"),
		DBUrl:      getEnv("DB_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
