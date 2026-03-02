package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort        string
	AppEnv            string
	DBUrl             string
	CORSAllowedOrigin string
	CORSAllowedMethods string
	CORSAllowedHeaders string
	RateLimitRPS      int
	RateLimitBurst    int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Env not found")
	}

	return &Config{
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		AppEnv:             getEnv("APP_ENV", "development"),
		DBUrl:              getEnv("DB_URL", ""),
		CORSAllowedOrigin:  getEnv("CORS_ALLOWED_ORIGIN", "*"),
		CORSAllowedMethods: getEnv("CORS_ALLOWED_METHODS", "GET, POST, PUT, DELETE, OPTIONS"),
		CORSAllowedHeaders: getEnv("CORS_ALLOWED_HEADERS", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization"),
		RateLimitRPS:       getEnvInt("RATE_LIMIT_RPS", 100),
		RateLimitBurst:     getEnvInt("RATE_LIMIT_BURST", 50),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
	}
	return fallback
}
