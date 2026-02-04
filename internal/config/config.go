package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBDatabase string
	ServerPort string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("failed to load env file: %w", err)
	}

	dbHost, err := getFromEnv("DB_HOST")
	if err != nil {
		return Config{}, err
	}

	dbPort, err := getFromEnv("DB_PORT")
	if err != nil {
		return Config{}, err
	}

	dbUser, err := getFromEnv("DB_USER")
	if err != nil {
		return Config{}, err
	}

	dbPassword, err := getFromEnv("DB_PASSWORD")
	if err != nil {
		return Config{}, err
	}

	dbDatabase, err := getFromEnv("DB_DATABASE")
	if err != nil {
		return Config{}, err
	}

	serverport, err := getFromEnv("SERVER_PORT")
	if err != nil {
		return Config{}, err
	}

	return Config{
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBDatabase: dbDatabase,
		ServerPort: serverport,
	}, nil
}

func getFromEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("environment variable %s not set", key)
	}
	return value, nil
}
