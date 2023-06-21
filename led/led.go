package led

import (
	"errors"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/whgentry/gomidi-led/animation"
	"github.com/whgentry/gomidi-led/control"
)

var (
	ErrLEDOutOfBounds = errors.New("led indicies are out of bounds")
	rowCount          int
	colCount          int
	frameRate         int

	Displays = []control.ProcessInterface[animation.PixelStateFrame, any]{
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
