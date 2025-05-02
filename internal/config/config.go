package config

import (
	"log"

	"github.com/joho/godotenv"
)

func Load() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, relying on env vars")
	}
}
