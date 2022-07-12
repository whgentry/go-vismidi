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
	FlushFrame() error
}

func HandleRefresh(ctx context.Context, lg LEDGridInterface, refreshRate int, frameFunc func(LEDGridInterface)) {
	frameDurationMs := 1000 / refreshRate
	ticker := time.NewTicker(time.Duration(frameDurationMs) * time.Millisecond)
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
