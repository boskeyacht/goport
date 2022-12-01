package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	RPC_URL     string
	DB_NAME     string
	DB_PASSWORD string
	HOST_PORT   string
	HOST_NAME   string
)

func getEnvOrFail(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing environment variable: %s", key)
	}
	return val
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	RPC_URL = getEnvOrFail("RPC_URL")
	DB_NAME = getEnvOrFail("DB_NAME")
	DB_PASSWORD = getEnvOrFail("DB_PASSWORD")
	HOST_PORT = getEnvOrFail("HOST_PORT")
	HOST_NAME = getEnvOrFail("HOST_NAME")
}
