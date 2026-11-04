package main

import (
	"log"
	"os"

	"pakku-agent/assets"
	"pakku-agent/internal/agent"
	"pakku-agent/internal/chat"
	"pakku-agent/internal/config"
	"pakku-agent/internal/openrouter"
	"pakku-agent/internal/tools"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := openrouter.NewClient(cfg.APIURL, cfg.APIKey, cfg.Model)
	toolManager := tools.NewManager(tools.PermRead, tools.PermWrite)
	pakku := agent.New(client, toolManager, cfg.MaxSteps, chat.Terminal{})

	assets.PrintBanner()
	return chat.Run(os.Stdin, pakku)
}
