package animations

import (
	"context"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/whgentry/gomidi-led/keyboard"
	"github.com/whgentry/gomidi-led/leds"
)

type PixelState struct {
	Color     colorful.Color
	Intensity float64
}

type Animation interface {
	Run(ctx context.Context)
}

var kboard *keyboard.Keyboard
var animations map[string]Animation
var activeAnimationName string
var numRows int
var numCols int
var pixels [][]*PixelState
var cancelActive func() = nil
var ctx context.Context = nil
var frameDuration time.Duration

func (ps *PixelState) Clear() {
	ps.Color = colorful.Color{0, 0, 0}
	ps.Intensity = 0
}

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

	activeAnimationName = name
	ctx, cancelActive = context.WithCancel(context.Background())
	go animations[activeAnimationName].Run(ctx)
}

func RotateAnimation() {
	Stop()
	if activeAnimationName == "velocity-bar" {
		SetAnimation("flowing-notes")
	} else {
		SetAnimation("velocity-bar")
	}
}

func Stop() {
	if ctx != nil {
		cancelActive()
		ctx = nil
		for i := range pixels {
			for _, ps := range pixels[i] {
				ps.Clear()
			}
		}
	}
}
