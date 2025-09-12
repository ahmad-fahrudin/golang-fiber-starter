package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret  string
}

var AppConfig *Config

func LoadConfig() error {
	// Load .env file if exists
	err := godotenv.Load()
	if err != nil {
		// If .env doesn't exist, continue with environment variables
	}

	AppConfig = &Config{
		Port:       getEnv("PORT", "3000"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "fiber_starter"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetPort() int {
	port, err := strconv.Atoi(AppConfig.Port)
	if err != nil {
		return 3000
	}
	return port
}

func GetDatabaseURL() string {
	return "host=" + AppConfig.DBHost +
		" port=" + AppConfig.DBPort +
		" user=" + AppConfig.DBUser +
		" password=" + AppConfig.DBPassword +
		" dbname=" + AppConfig.DBName +
		" sslmode=" + AppConfig.DBSSLMode
}
