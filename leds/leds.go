package leds

import (
	"errors"
)

var (
	ErrLEDOutOfBounds = errors.New("led indicies are out of bounds")
)

type LEDGridInterface interface {
	SetLED(row int, col int, color Color) error
	UpdateLEDs() error
}
