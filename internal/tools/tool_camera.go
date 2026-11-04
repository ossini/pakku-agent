package tools

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"
)

func init() {
	register(Tool{
		Name:        "take_photo",
		Title:       "Foto aufnehmen",
		Description: "Takes a photo with the Raspberry Pi camera. The image is attached to the conversation right after this tool result.",
		Permission:  PermRead,
		Parameters: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
		Run: takePhoto,
	})
}

func takePhoto(ctx context.Context, _ string) (Result, error) {
	cmd := exec.CommandContext(ctx, "rpicam-still", "-n", "-t", "1000", "-e", "jpg", "--width", "1280", "--height", "960", "-o", "-")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return Result{}, fmt.Errorf("rpicam-still: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return Result{
		Text:     fmt.Sprintf("Photo taken (%d KB).", stdout.Len()/1024),
		ImageURL: "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(stdout.Bytes()),
	}, nil
}
