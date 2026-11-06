package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	apiURL           = "https://openrouter.ai/api/v1/responses"
	defaultMaxSteps  = 20
	defaultRelayChip = "gpiochip0"
)

type Config struct {
	APIURL    string
	APIKey    string
	Model     string
	MaxSteps  int
	RelayChip string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	model := os.Getenv("MODEL")
	if apiKey == "" || model == "" {
		return Config{}, errors.New("API_KEY or MODEL is missing in your .env")
	}

	maxSteps, err := strconv.Atoi(os.Getenv("MAX_STEPS"))
	if err != nil || maxSteps <= 0 {
		maxSteps = defaultMaxSteps
	}
	relayChip := os.Getenv("RELAY_GPIO_CHIP")
	if relayChip == "" {
		relayChip = defaultRelayChip
	}

	return Config{
		APIURL:    apiURL,
		APIKey:    apiKey,
		Model:     model,
		MaxSteps:  maxSteps,
		RelayChip: relayChip,
	}, nil
}
