package leds

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/nsf/termbox-go"
)

type VirtualLEDGrid struct {
	NumRows int
	NumCols int
	Grid    [][]*LED
}

func NewVirtualLEDGrid(numRows int, numCols int) *VirtualLEDGrid {
	ledGrid := &VirtualLEDGrid{
		NumRows: numRows,
		NumCols: numCols,
		Grid:    make([][]*LED, numRows),
	}
	for i := range ledGrid.Grid {
		ledGrid.Grid[i] = make([]*LED, numCols)
		for j := range ledGrid.Grid[i] {
			ledGrid.Grid[i][j] = &LED{}
		}
	}
	return ledGrid
}

func (lg *VirtualLEDGrid) SetLED(row int, col int, color colorful.Color) error {
	if row < 0 || row > lg.NumRows || col < 0 || col > lg.NumCols {
		return ErrLEDOutOfBounds
	}
	lg.Grid[row][col].Color = color
	return nil
}

func (lg *VirtualLEDGrid) FlushFrame() error {
	termbox.SetOutputMode(termbox.OutputRGB)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for i := range lg.Grid {
		for j, led := range lg.Grid[i] {
			if !IsColorOff(led.Color) {
				r, g, b := led.Color.RGB255()
				fg := termbox.RGBToAttribute(r, g, b)
				bg := termbox.Attribute(termbox.ColorDefault)
				termbox.SetCell(j, i, '*', fg, bg)
			}
		}
	}
	termbox.Flush()
	return nil
}

func (lg *VirtualLEDGrid) ClearFrame() error {
	for i := range lg.Grid {
		for _, led := range lg.Grid[i] {
			led.Color = ColorOff()
		}
	}
	return nil
}
