package tools

import (
	"context"
	"testing"
)

func TestNewManagerFiltersByPermission(t *testing.T) {
	if m := NewManager(PermDangerous); m.Title("get_time") != "get_time" {
		t.Error("get_time must not be available without read permission")
	}
	if m := NewManager(PermRead); m.Title("get_time") != "Aktuelle Uhrzeit" {
		t.Error("get_time must be available with read permission")
	}
}

func TestRunRejectsUnknownTool(t *testing.T) {
	if _, err := NewManager(PermRead).Run(context.Background(), "nope", ""); err == nil {
		t.Error("expected error for unknown tool")
	}
}

func TestGetTimeRejectsInvalidTimezone(t *testing.T) {
	if _, err := getTime(context.Background(), `{"timezone":"Mars/Olympus"}`); err == nil {
		t.Error("expected error for invalid timezone")
	}
}
