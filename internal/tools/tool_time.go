package tools

import (
	"context"
	"encoding/json"
	"time"
)

func init() {
	register(Tool{
		Name:        "get_time",
		Title:       "Aktuelle Uhrzeit",
		Description: "Returns the current date and time. Optionally in a given IANA timezone like Europe/Berlin.",
		Permission:  PermRead,
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"timezone": map[string]any{
					"type":        "string",
					"description": "IANA timezone, e.g. Europe/Berlin",
				},
			},
		},
		Run: getTime,
	})
}

func getTime(_ context.Context, args string) (Result, error) {
	var in struct {
		Timezone string `json:"timezone"`
	}
	if args != "" {
		if err := json.Unmarshal([]byte(args), &in); err != nil {
			return Result{}, err
		}
	}

	now := time.Now()
	if in.Timezone != "" {
		loc, err := time.LoadLocation(in.Timezone)
		if err != nil {
			return Result{}, err
		}
		now = now.In(loc)
	}
	return Result{Text: now.Format("02.01.2006 15:04:05 MST")}, nil
}
