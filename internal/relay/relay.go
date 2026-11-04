package relay

import (
	"errors"
	"fmt"
	"sync"
)

type lines interface {
	SetValues([]int) error
	Close() error
}

// Controller controls the three Waveshare relay channels (1 through 3).
type Controller struct {
	mu     sync.Mutex
	lines  lines
	values [3]int
	closed bool
}

func (c *Controller) Set(channel int, on bool) error {
	if channel < 1 || channel > 3 {
		return fmt.Errorf("invalid relay channel %d: expected 1, 2 or 3", channel)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return errors.New("relay controller is closed")
	}
	values := c.values
	if on {
		values[channel-1] = 1
	} else {
		values[channel-1] = 0
	}
	if err := c.lines.SetValues(values[:]); err != nil {
		return fmt.Errorf("set relay channel %d: %w", channel, err)
	}
	c.values = values
	return nil
}

// Close requests all channels off before releasing the GPIO lines.
func (c *Controller) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	return errors.Join(c.lines.SetValues([]int{0, 0, 0}), c.lines.Close())
}
