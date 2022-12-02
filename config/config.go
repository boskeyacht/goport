package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	RPC_URL   string
	DB_NAME   string
	HOST_PORT string
	HOST_NAME string
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
		log.Printf("Error loading .env file: %v", err.Error())
	}

	RPC_URL = getEnvOrFail("RPC_URL")
	DB_NAME = getEnvOrFail("DB_NAME")
	HOST_PORT = getEnvOrFail("HOST_PORT")
	HOST_NAME = getEnvOrFail("HOST_NAME")
}
