package leds

import (
	"github.com/nsf/termbox-go"
)

type VirtualLED struct {
	Color Color
}

type VirtualLEDGrid struct {
	NumRows int
	NumCols int
	Grid    [][]*VirtualLED
}

func NewVirtualLEDGrid(numRows int, numCols int) *VirtualLEDGrid {
	ledGrid := &VirtualLEDGrid{
		NumRows: numRows,
		NumCols: numCols,
		Grid:    make([][]*VirtualLED, numRows),
	}
	for i := range ledGrid.Grid {
		ledGrid.Grid[i] = make([]*VirtualLED, numCols)
		for j := range ledGrid.Grid[i] {
			ledGrid.Grid[i][j] = &VirtualLED{}
		}
	}
	return ledGrid
}

func (lg *VirtualLEDGrid) SetLED(row int, col int, color Color) error {
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
			if led.Color != Off {
				// fg := termbox.Attribute(led.Color)
				fg := termbox.RGBToAttribute(
					uint8(led.Color>>16&0xFF),
					uint8(led.Color>>8&0xFF),
					uint8(led.Color&0xFF))
				bg := termbox.Attribute(Off)
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
			led.Color = Off
		}
	}
	return nil
}
