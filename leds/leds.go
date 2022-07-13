package leds

import (
	"context"
	"errors"
	"time"
)

var (
	ErrLEDOutOfBounds = errors.New("led indicies are out of bounds")
)

type LEDGridInterface interface {
	SetLED(row int, col int, color Color) error
	ClearFrame() error
	FlushFrame() error
}

func HandleRefresh(ctx context.Context, lg LEDGridInterface, refreshRate int, frameFunc func(LEDGridInterface)) {
	ticker := time.NewTicker(time.Second / time.Duration(refreshRate))
	for {
		select {
		case <-ticker.C:
			frameFunc(lg)
			lg.FlushFrame()
		case <-ctx.Done():
			return
		}
	}
}
