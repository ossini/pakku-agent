package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	APIUrl string
	APIKey string
	Model string
	MaxSteps int
}

func Load() Config {
	apiUrl := "https://openrouter.ai/api/v1/responses"	
	godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	model := os.Getenv("MODEL")
	if apiKey == "" || model == "" {
		log.Fatal("API_KEY or MODEL is missing in your .env")
	}

	maxSteps, err := strconv.Atoi(os.Getenv("MAX_STEPS"))
	if err != nil || maxSteps <= 0 {
		maxSteps = 20
	}

	return Config{
		APIUrl: apiUrl,
		APIKey:   apiKey,
		Model:    model,
		MaxSteps: maxSteps,
	}
}
