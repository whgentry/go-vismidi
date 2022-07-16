package animations

import (
	"context"
	"sync"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/whgentry/gomidi-led/keyboard"
	"github.com/whgentry/gomidi-led/leds"
)

type RunFunc func(ctx context.Context, settings Settings)

type PixelState struct {
	Color     colorful.Color
	Intensity float64
}

type Settings struct {
	LowerColor colorful.Color
	UpperColor colorful.Color
}

type Animation struct {
	Name        string
	Key         string
	Description string
	Settings    Settings
	Run         RunFunc
}

var (
	kboard        *keyboard.Keyboard
	numRows       int
	numCols       int
	pixels        [][]*PixelState
	frameDuration time.Duration

	cancelActive  func()          = func() {}
	ctxActive     context.Context = nil
	cancelControl func()          = func() {}
	ctxControl    context.Context = nil

	wg                   sync.WaitGroup
	activeAnimationIndex int
	activeAnimationChan  chan int
	stopAnimation        chan bool

	CommonSettings Settings
)

var animations = []*Animation{
	VelocityBar,
	VelocityBarMirror,
	FlowingNotes,
}

func (ps *PixelState) Clear() {
	ps.Color = colorful.Color{R: 0, G: 0, B: 0}
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

	// Initialize Settings
	CommonSettings = Settings{
		LowerColor: colorful.FastLinearRgb(0, 1, 0),
		UpperColor: colorful.FastLinearRgb(1, 0, 0),
	}
	for _, animation := range animations {
		animation.Settings = CommonSettings
	}

	activeAnimationChan = make(chan int)
	ctxControl, cancelControl = context.WithCancel(context.Background())
	go HandleAnimationControl(ctxControl)
}

func Close() {
	cancelControl()
}

func HandleAnimationControl(ctx context.Context) {
	stop := func() {
		if ctxActive != nil {
			wg.Add(1)
			cancelActive()
			wg.Wait()
			ctxActive = nil
			for i := range pixels {
				for _, ps := range pixels[i] {
					ps.Clear()
				}
			}
		}
	}
	for {
		select {
		case <-stopAnimation:
			stop()
		case activeAnimationIndex = <-activeAnimationChan:
			stop()
			ctxActive, cancelActive = context.WithCancel(ctxControl)
			go animations[activeAnimationIndex].Run(ctxActive, animations[activeAnimationIndex].Settings)
		case <-ctx.Done():
			return
		}
	}
}

func StopAnimation() {
	stopAnimation <- true
}

func FrameHandler(lg leds.LEDGridInterface) {
	for i := range pixels {
		for j := range pixels[i] {
			lg.SetLED(i, j, pixels[i][j].Color)
		}
	}
}

func SetAnimationByIndex(index int) {
	if index >= 0 && index < len(animations) {
		activeAnimationChan <- index
	}
}

func SetAnimationByName(name string) {
	for i, animation := range animations {
		if animation.Key == name {
			activeAnimationChan <- i
			break
		}
	}
}

func PreviousAnimation() {
	activeAnimationChan <- (activeAnimationIndex + len(animations) - 1) % len(animations)
}

func NextAnimation() {
	activeAnimationChan <- (activeAnimationIndex + 1) % len(animations)
}
