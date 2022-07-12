package animations

import (
	"context"

	"github.com/whgentry/gomidi-led/keyboard"
	"github.com/whgentry/gomidi-led/leds"
)

type PixelState struct {
	Color     leds.Color
	Intensity float32
}

type Animation interface {
	FrameHandler(lg leds.LEDGridInterface)
	Run(ctx context.Context)
}

var kboard *keyboard.Keyboard
var animations map[string]Animation
var numRows int
var numCols int
var pixels [][]*PixelState

func Initialize(rows int, cols int, k *keyboard.Keyboard) {
	numRows = rows
	numCols = cols
	kboard = k
	pixels = make([][]*PixelState, numRows)
	for i := range pixels {
		pixels[i] = make([]*PixelState, numCols)
		for j := range pixels[i] {
			pixels[i][j] = &PixelState{}
		}
	}

	animations = map[string]Animation{
		"velocity-bar":  &VelocityBarAnimation{},
		"flowing-notes": &FlowingNotesAnimation{},
	}
}

func GetAnimation(name string) Animation {
	return animations[name]
}
