package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config хранит все настройки для приложения
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

// Load загружает настройки из env
func Load() *Config {
	// Загружаем env. Если файла нет (например на продакшене),
	// ошибка не критична, так как переменные могут быть заданы через систему\
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	return &Config{
		Port:        getEnv("PORT", ":8080"),
		DatabaseURL: getEnv("DB_URL", ""),
		JWTSecret:   getEnv("JWT_Secret", "secret"),
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return fallback
}
