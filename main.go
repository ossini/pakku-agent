package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	APIKey := os.Getenv("OPENROUTER_API_KEY")

}