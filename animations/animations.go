package animations

import (
	"context"
	"time"

	"github.com/whgentry/gomidi-led/keyboard"
	"github.com/whgentry/gomidi-led/leds"
)

type PixelState struct {
	Color     leds.Color
	Intensity float32
}

type Animation interface {
	Run(ctx context.Context)
}

var kboard *keyboard.Keyboard
var animations map[string]Animation
var active Animation
var numRows int
var numCols int
var pixels [][]*PixelState
var cancelActive func() = nil
var ctx context.Context = nil
var frameDuration time.Duration

func Initialize(rows int, cols int, rate int, k *keyboard.Keyboard) {
	frameDuration = time.Second / time.Duration(rate)
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

func FrameHandler(lg leds.LEDGridInterface) {
	for i := range pixels {
		for j := range pixels[i] {
			lg.SetLED(i, j, pixels[i][j].Color)
		}
	}
}

func SetAnimation(name string) {
	if ctx != nil {
		cancelActive()
	}

	active = animations[name]
	ctx, cancelActive = context.WithCancel(context.Background())
	go active.Run(ctx)
}

func Stop() {
	if ctx != nil {
		cancelActive()
		ctx = nil
	}
}
