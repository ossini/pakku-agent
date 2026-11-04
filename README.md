# pakku-agent
A smart plant agent with openrouter integration

## Setup

```sh
cp .env.example .env
go run .
```

| Variable    | Description                                  |
|-------------|----------------------------------------------|
| `API_KEY`   | OpenRouter API key                           |
| `MODEL`     | OpenRouter model with tool support           |
| `MAX_STEPS` | Model requests per message, defaults to `20` |

## Adding a tool

Create `internal/tools/tool_<name>.go` and register the tool in `init`:

```go
package tools

import "context"

func init() {
	register(Tool{
		Name:        "get_plant",
		Title:       "Pflanze abrufen",
		Description: "Returns the current state of the plant.",
		Permission:  PermRead,
		Parameters: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
		Run: getPlant,
	})
}

func getPlant(ctx context.Context, args string) (Result, error) {
	return Result{Text: "healthy"}, nil
}
```

`args` is the JSON object the model sent. `Title` is shown in the chat, the
other fields are sent to the model. Returned errors are shown in the chat and
passed to the model. Set `Result.ImageURL` to a data URL to send an image to
the model.

## Camera

`take_photo` calls `rpicam-still`, so it needs a Raspberry Pi with a connected
camera and a `MODEL` with image input.

`main.go` enables `PermRead` and `PermWrite`. `PermDangerous` tools stay
unavailable until they are added there.
