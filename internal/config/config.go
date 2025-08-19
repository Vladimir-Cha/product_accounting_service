package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	DBUrl      string
	ServerPort int
	DBMaxConns int32
	DBMinConns int32
	DBConnLife time.Duration
}

func Load() (*DBConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	// вспомогательные переменные для формирования url
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		return nil, fmt.Errorf("missing parametrs: DB_USER, DB_PASSWORD, DB_HOST, DB_PORT or DB_NAME")
	}

	// формируем url из переменных окружения
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	cfg := &DBConfig{
		DBUrl:      dbURL,
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
		DBMaxConns: getEnvAsInt32("DB_MAX_CONNS", 10),
		DBMinConns: getEnvAsInt32("DB_MIN_CONNS", 2),
		DBConnLife: getEnvAsDuration("DB_CONN_LIFE", time.Hour),
	}

	return cfg, nil
}

// Вспомогательные функции для парсинга
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsInt32(key string, defaultValue int32) int32 {
	return int32(getEnvAsInt(key, int(defaultValue)))
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
