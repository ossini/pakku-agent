package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"pakku-agent/internal/openrouter"
	"pakku-agent/internal/tools"
)

const instructions = "You are Pakku, a smart plant agent. Use the available tools whenever they help, and answer in the language of the user. Reply in plain text without Markdown, because the chat runs in a terminal."

type Observer interface {
	Thinking()
	ToolStarted(title, args string)
	ToolFinished(output string, err error)
}

type Agent struct {
	client   *openrouter.Client
	tools    *tools.Manager
	maxSteps int
	observer Observer
	history  []json.RawMessage
}

func New(client *openrouter.Client, tools *tools.Manager, maxSteps int, observer Observer) *Agent {
	return &Agent{client: client, tools: tools, maxSteps: maxSteps, observer: observer}
}

// Send runs one conversation turn. The history is only updated when the turn
// succeeds, so a failed or cancelled turn leaves no unanswered tool calls behind.
func (a *Agent) Send(ctx context.Context, text string) (string, error) {
	history := append(a.history, openrouter.UserMessage(text))

	for range a.maxSteps {
		a.observer.Thinking()
		resp, err := a.client.Respond(ctx, openrouter.Request{
			Instructions: instructions,
			Input:        history,
			Tools:        a.tools.Definitions(),
		})
		if err != nil {
			return "", err
		}

		history = append(history, resp.Items...)
		if len(resp.ToolCalls) == 0 {
			a.history = history
			return resp.Text, nil
		}

		var images []json.RawMessage
		for _, call := range resp.ToolCalls {
			if err := ctx.Err(); err != nil {
				return "", err
			}
			result := a.runTool(ctx, call)
			history = append(history, openrouter.ToolOutput(call.ID, result.Text))
			if result.ImageURL != "" {
				images = append(images, openrouter.UserImage(result.ImageURL))
			}
		}
		history = append(history, images...)
	}
	return "", fmt.Errorf("no final answer after %d steps", a.maxSteps)
}

func (a *Agent) runTool(ctx context.Context, call openrouter.ToolCall) tools.Result {
	a.observer.ToolStarted(a.tools.Title(call.Name), call.Arguments)
	result, err := a.tools.Run(ctx, call.Name, call.Arguments)
	a.observer.ToolFinished(result.Text, err)
	if err != nil {
		return tools.Result{Text: "error: " + err.Error()}
	}
	return result
}
