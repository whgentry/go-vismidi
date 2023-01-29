package leds

import (
	"errors"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/whgentry/gomidi-led/animations"
	"github.com/whgentry/gomidi-led/control"
)

// type LED struct {
// 	Color colorful.Color
// }

// type LEDGridInterface interface {
// 	SetLED(row int, col int, color colorful.Color) error
// 	ClearFrame() error
// 	FlushFrame() error
// }

var (
	ErrLEDOutOfBounds = errors.New("led indicies are out of bounds")
	rowCount          int
	colCount          int
	frameRate         int

	Displays = []control.ProcessInterface[animations.PixelStateFrame, any]{
		&VirtualLEDGrid{},
	}
)

func Initialize(numRows int, numCols int, rate int) {
	rowCount = numRows
	colCount = numCols
	frameRate = rate
}

func ColorOff() colorful.Color {
	color, _ := colorful.Hex("#000")
	return color
}

func IsColorOff(color colorful.Color) bool {
	r, g, b, _ := color.RGBA()
	return r == 0 && g == 0 && b == 0
}

// func HandleRefresh(ctx context.Context, lg LEDGridInterface, refreshRate int, frameFunc func(LEDGridInterface)) {
// 	ticker := time.NewTicker(time.Second / time.Duration(refreshRate))
// 	for {
// 		select {
// 		case <-ticker.C:
// 			frameFunc(lg)
// 			lg.FlushFrame()
// 		case <-ctx.Done():
// 			return
// 		}
// 	}
// }
