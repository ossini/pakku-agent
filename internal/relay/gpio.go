package relay

import (
	"fmt"

	"github.com/warthog618/go-gpiocdev"
)

// New requests BCM GPIO 26, 20 and 21 together, initially inactive.
// Active-low translates logical on/off to the board's electrical levels.
func New(chip string) (*Controller, error) {
	if chip == "" {
		return nil, fmt.Errorf("GPIO chip is required")
	}
	ll, err := gpiocdev.RequestLines(chip, []int{26, 20, 21},
		gpiocdev.AsActiveLow, gpiocdev.AsOutput(0, 0, 0), gpiocdev.WithConsumer("pakku-agent"))
	if err != nil {
		return nil, fmt.Errorf("request relay GPIOs on %s: %w", chip, err)
	}
	return &Controller{lines: ll}, nil
}
