package tools

import (
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

func getTime(args string) (string, error) {
	var in struct {
		Timezone string `json:"timezone"`
	}
	if args != "" {
		if err := json.Unmarshal([]byte(args), &in); err != nil {
			return "", err
		}
	}

	now := time.Now()
	if in.Timezone != "" {
		loc, err := time.LoadLocation(in.Timezone)
		if err != nil {
			return "", err
		}
		now = now.In(loc)
	}
	return now.Format("02.01.2006 15:04:05 MST"), nil
}
