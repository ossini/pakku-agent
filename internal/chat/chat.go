package chat

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/fatih/color"

	"pakku-agent/internal/agent"
)

const (
	maxLineSize    = 1 << 20
	maxPreviewSize = 200
	clearLine      = "\r\033[K"
)

var (
	userStyle  = color.New(color.FgGreen, color.Bold)
	agentStyle = color.New(color.FgYellow, color.Bold)
	toolStyle  = color.New(color.FgCyan)
	mutedStyle = color.New(color.Faint)
	errorStyle = color.New(color.FgRed)
)

func Run(in io.Reader, pakku *agent.Agent) error {
	mutedStyle.Println("Beenden mit exit oder Strg+D, laufende Antwort abbrechen mit Strg+C.")

	scanner := bufio.NewScanner(in)
	scanner.Buffer(nil, maxLineSize)
	for {
		userStyle.Print("du › ")
		if !scanner.Scan() {
			fmt.Println()
			return scanner.Err()
		}

		text := strings.TrimSpace(scanner.Text())
		switch text {
		case "":
			continue
		case "exit", "quit":
			return nil
		}
		send(pakku, text)
	}
}

func send(pakku *agent.Agent, text string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	reply, err := pakku.Send(ctx, text)
	clearStatus()

	switch {
	case errors.Is(err, context.Canceled):
		mutedStyle.Println("Abgebrochen.")
	case err != nil:
		errorStyle.Printf("Fehler: %v\n", err)
	default:
		agentStyle.Print("pakku › ")
		fmt.Println(reply)
	}
}

type Terminal struct{}

func (Terminal) Thinking() {
	if !color.NoColor {
		mutedStyle.Print("denkt nach…")
	}
}

func (Terminal) ToolStarted(title, args string) {
	clearStatus()
	if args == "" || args == "{}" {
		toolStyle.Printf("⚙ %s\n", title)
		return
	}
	toolStyle.Printf("⚙ %s %s\n", title, preview(args))
}

func (Terminal) ToolFinished(output string, err error) {
	if err != nil {
		errorStyle.Printf("  ✗ %s\n", preview(err.Error()))
		return
	}
	mutedStyle.Printf("  ↳ %s\n", preview(output))
}

func clearStatus() {
	if !color.NoColor {
		fmt.Print(clearLine)
	}
}

func preview(s string) string {
	line := []rune(strings.Join(strings.Fields(s), " "))
	if len(line) <= maxPreviewSize {
		return string(line)
	}
	return string(line[:maxPreviewSize]) + "…"
}
