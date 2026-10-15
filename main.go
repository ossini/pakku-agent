package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"pakku-agent/assets"
	"pakku-agent/internal/config"
)

func main() {
	cfg := config.Load()

	body, err := json.Marshal(map[string]string{
		"input": "Tell me a joke",
		"model": cfg.Model,
	})
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", cfg.APIUrl, bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

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