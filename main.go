package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"pakku-agent/assets"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	const openRouterUrl = "https://openrouter.ai/api/v1/responses"
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	model := os.Getenv("OPENROUTER_MODEL")

	if apiKey == "" || model == "" {
		log.Fatal("OPENROUTER_API_KEY or OPENROUTER_MODEL is missing in your .env")
	}

	body, err := json.Marshal(map[string]string{
		"input": "Tell me a joke",
		"model": model,
	})
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", openRouterUrl, bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Fehler %d: %s", resp.StatusCode, respBody)
	}

	var out bytes.Buffer
	if err := json.Indent(&out, respBody, "", "  "); err != nil {
		fmt.Println(string(respBody))
		return
	}
	assets.PrintBanner()
	fmt.Println(out.String())
}