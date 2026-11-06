# pakku-agent
A smart plant agent with openrouter integration

## Setup

```sh
cp .env.example .env
go run .
```

| Variable          | Description                                                    |
|-------------------|----------------------------------------------------------------|
| `API_KEY`         | OpenRouter API key                                             |
| `MODEL`           | OpenRouter model with tool support                             |
| `MAX_STEPS`       | Model requests per message, defaults to `20`                   |
| `RELAY_GPIO_CHIP` | GPIO chip with the Pi header pins; defaults to `gpiochip0` |

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

## Relay board

The default `RELAY_GPIO_CHIP=gpiochip0` targets the Pi 5 used for this project
(`pinctrl-rp1`). Check that BCM GPIO 26, 20 and 21 belong to that chip on your
Pi; chip numbers can change with the OS or kernel. Override `RELAY_GPIO_CHIP`
in `.env` if needed (`gpiochip4` is a symlink to `gpiochip0` on the Pi 5).
The `set_relay` tool is registered at startup. It accepts
`{"channel":1,"on":true}` (channels 1-3); `on=false` switches it off. The
Waveshare RPi Relay Board uses BCM 26/20/21 for CH1/CH2/CH3 and is active LOW.
The GPIO device must be accessible to the process. The program starts with all
three relays off and switches them off again on normal exit (`exit` or Ctrl+D).
The relay driver is Linux-only and requires the GPIO character device API
(kernel 5.10+); non-Linux builds are not supported. Check
the board's `Relay_JMP` jumper, wiring and connected load before enabling it.
