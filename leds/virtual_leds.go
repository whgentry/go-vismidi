package leds

import (
	"context"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/whgentry/gomidi-led/animations"
)

type VirtualLEDGrid struct {
	state animations.PixelStateFrame
}

func (lg *VirtualLEDGrid) Run(ctx context.Context, input chan animations.PixelStateFrame, out chan any) {
	ticker := time.NewTicker(time.Second / time.Duration(frameRate))
	for {
		select {
		case state := <-input:
			lg.state = state
		case <-ticker.C:
			for i := range lg.state.Pixels {
				for j, led := range lg.state.Pixels[i] {
					if !IsColorOff(led.Color) {
						r, g, b := led.Color.RGB255()
						fg := termbox.RGBToAttribute(r, g, b)
						bg := termbox.Attribute(termbox.ColorDefault)
						_, height := termbox.Size()
						termbox.SetCell(j*2, height-i, '*', fg, bg)
						termbox.SetCell(j*2+1, height-i, '*', fg, bg)
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (lg *VirtualLEDGrid) Name() string {
	return "VirtualLEDGrid"
}

// func NewVirtualLEDGrid(numRows int, numCols int) *VirtualLEDGrid {
// 	ledGrid := &VirtualLEDGrid{
// 		NumRows: numRows,
// 		NumCols: numCols,
// 		Grid:    make([][]*LED, numRows),
// 	}
// 	for i := range ledGrid.Grid {
// 		ledGrid.Grid[i] = make([]*LED, numCols)
// 		for j := range ledGrid.Grid[i] {
// 			ledGrid.Grid[i][j] = &LED{}
// 		}
// 	}
// 	return ledGrid
// }

// func (lg *VirtualLEDGrid) SetLED(row int, col int, color colorful.Color) error {
// 	if row < 0 || row > lg.NumRows || col < 0 || col > lg.NumCols {
// 		return ErrLEDOutOfBounds
// 	}
// 	lg.Grid[row][col].Color = color
// 	return nil
// }

// func (lg *VirtualLEDGrid) FlushFrame() error {
// 	termbox.SetOutputMode(termbox.OutputRGB)
// 	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
// 	for i := range lg.Grid {
// 		for j, led := range lg.Grid[i] {
// 			if !IsColorOff(led.Color) {
// 				r, g, b := led.Color.RGB255()
// 				fg := termbox.RGBToAttribute(r, g, b)
// 				bg := termbox.Attribute(termbox.ColorDefault)
// 				_, height := termbox.Size()
// 				termbox.SetCell(j*2, height-i, '*', fg, bg)
// 				termbox.SetCell(j*2+1, height-i, '*', fg, bg)
// 			}
// 		}
// 	}
// 	termbox.Flush()
// 	return nil
// }

// func (lg *VirtualLEDGrid) ClearFrame() error {
// 	for i := range lg.st {
// 		for _, led := range lg.Grid[i] {
// 			led.Color = ColorOff()
// 		}
// 	}
// 	return nil
// }
