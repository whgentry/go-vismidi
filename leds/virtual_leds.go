package leds

import (
	"context"

	"github.com/ansoni/termination"
)

type VirtualLEDGrid struct {
	NumRows int
	NumCols int
	Grid    [][]Color
	Term    *termination.Termination
}

var ledShape = termination.Shape{
	"default": []string{""},
	"w":       []string{"*"},
	"r":       []string{"*"},
	"g":       []string{"*"},
	"b":       []string{"*"},
}

var ledColorMask = map[string][]string{
	"default": {"w"},
	"w":       {"w"},
	"r":       {"r"},
	"g":       {"g"},
	"b":       {"b"},
}

func ledMovement(t *termination.Termination, e *termination.Entity, position termination.Position) termination.Position {
	lg := e.Data.(*VirtualLEDGrid)
	color := lg.Grid[position.Y][position.X]

	switch color {
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
		for j := range lg.Grid[i] {
			ledEntity := lg.Term.NewEntity(termination.Position{X: j, Y: i, Z: 0})
			ledEntity.Shape = ledShape
			ledEntity.ColorMask = ledColorMask
			ledEntity.MovementCallback = ledMovement
			ledEntity.Data = lg
		}
	}
	go lg.Term.Animate()

	for {
		select {
		case <-ctx.Done():
			lg.Term.Close()
		}
	}
}

func NewVirtualLEDGrid(numRows int, numCols int) *VirtualLEDGrid {
	ledGrid := &VirtualLEDGrid{
		NumRows: numRows,
		NumCols: numCols,
		Grid:    make([][]Color, numRows),
	}
	for i := range ledGrid.Grid {
		ledGrid.Grid[i] = make([]Color, numCols)
	}
	return ledGrid
}

func (lg *VirtualLEDGrid) SetLED(row int, col int, color Color) error {
	if row < 0 || row > lg.NumRows || col < 0 || col > lg.NumCols {
		return ErrLEDOutOfBounds
	}
	lg.Grid[row][col] = color
	return nil
}

func (lg *VirtualLEDGrid) FlushFrame() error {
	return nil
}
