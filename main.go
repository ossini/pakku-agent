package main

import (
	"errors"
	"log"
	"os"

	"pakku-agent/assets"
	"pakku-agent/internal/agent"
	"pakku-agent/internal/chat"
	"pakku-agent/internal/config"
	"pakku-agent/internal/openrouter"
	"pakku-agent/internal/relay"
	"pakku-agent/internal/tools"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := openrouter.NewClient(cfg.APIURL, cfg.APIKey, cfg.Model)
	toolManager := tools.NewManager(tools.PermRead, tools.PermWrite)
	if cfg.RelayChip != "" {
		controller, openErr := relay.New(cfg.RelayChip)
		if openErr != nil {
			return openErr
		}
		defer func() { err = errors.Join(err, controller.Close()) }()
		toolManager.ConfigureRelay(controller)
	}
	pakku := agent.New(client, toolManager, cfg.MaxSteps, chat.Terminal{})

	assets.PrintBanner()
	return chat.Run(os.Stdin, pakku)
}
