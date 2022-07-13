package leds

import (
	"context"

	"github.com/ansoni/termination"
)

type VirtualLED struct {
	Color Color
}

type VirtualLEDGrid struct {
	NumRows int
	NumCols int
	Grid    [][]*VirtualLED
	Term    *termination.Termination
}

var ledShape = termination.Shape{
	"default": []string{""},
	"w":       []string{"**"},
	"r":       []string{"**"},
	"g":       []string{"**"},
	"b":       []string{"**"},
}

var ledColorMask = map[string][]string{
	"default": {"ww"},
	"w":       {"ww"},
	"r":       {"rr"},
	"g":       {"gg"},
	"b":       {"bb"},
}

func ledMovement(t *termination.Termination, e *termination.Entity, position termination.Position) termination.Position {
	led := e.Data.(*VirtualLED)

	switch led.Color {
	case White:
		e.ShapePath = "w"
	case Red:
		e.ShapePath = "r"
	case Green:
		e.ShapePath = "g"
	case Blue:
		e.ShapePath = "b"
	case Off:
		e.ShapePath = "default"
	default:
		e.ShapePath = "default"
	}
	return position
}

func AnimateVirtualGrid(ctx context.Context, lg *VirtualLEDGrid, framesPerSecond int) {
	lg.Term = termination.New()
	lg.Term.FramesPerSecond = framesPerSecond
	for i := range lg.Grid {
		for j, led := range lg.Grid[i] {
			ledEntity := lg.Term.NewEntity(termination.Position{
				X: j * 2,
				Y: lg.Term.Height - i,
				Z: 0,
			})
			ledEntity.Shape = ledShape
			ledEntity.ColorMask = ledColorMask
			ledEntity.MovementCallback = ledMovement
			ledEntity.Data = led
		}
	}
	go lg.Term.Animate()
	<-ctx.Done()
	lg.Term.Close()
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
