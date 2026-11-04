package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type RelayController interface {
	Set(channel int, on bool) error
}

// RelayTool is registered only when relay hardware is enabled at startup.
func RelayTool(controller RelayController) Tool {
	return Tool{
		Name:        "set_relay",
		Title:       "Relais schalten",
		Description: "Switches a Waveshare RPi Relay Board channel. Channel is 1, 2 or 3, not a GPIO number. on=true energizes the relay; on=false releases it.",
		Permission:  PermWrite,
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"channel": map[string]any{"type": "integer", "enum": []int{1, 2, 3}},
				"on":      map[string]any{"type": "boolean"},
			},
			"required":             []string{"channel", "on"},
			"additionalProperties": false,
		},
		Run: func(ctx context.Context, args string) (Result, error) {
			var in struct {
				Channel *int  `json:"channel"`
				On      *bool `json:"on"`
			}
			if err := json.Unmarshal([]byte(args), &in); err != nil {
				return Result{}, fmt.Errorf("invalid relay arguments: %w", err)
			}
			if in.Channel == nil || in.On == nil {
				return Result{}, errors.New("channel and on are required")
			}
			if *in.Channel < 1 || *in.Channel > 3 {
				return Result{}, fmt.Errorf("invalid relay channel %d: expected 1, 2 or 3", *in.Channel)
			}
			if err := ctx.Err(); err != nil {
				return Result{}, err
			}
			if err := controller.Set(*in.Channel, *in.On); err != nil {
				return Result{}, err
			}
			state := "off"
			if *in.On {
				state = "on"
			}
			return Result{Text: fmt.Sprintf("Relay channel %d is %s.", *in.Channel, state)}, nil
		},
	}
}
