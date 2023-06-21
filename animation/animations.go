package animation

import (
	"context"

	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/whgentry/gomidi-led/control"
	"github.com/whgentry/gomidi-led/midi"
)

type PixelState struct {
	Color     colorful.Color
	Intensity float64
}

type PixelStateFrame struct {
	Pixels   [][]PixelState
	RowCount int
	ColCount int
}

type CommonSettings struct {
	LowerColor colorful.Color
	UpperColor colorful.Color
}

type Settings struct {
	*CommonSettings
}

type Animation struct {
	Description string
	Settings    *Settings
	name        string
	run         control.RunFunc[midi.MIDIEvent, PixelStateFrame]
}

func (a *Animation) Run(ctx context.Context, input chan midi.MIDIEvent, out chan PixelStateFrame) {
	a.run(ctx, input, out)
}

func (a *Animation) Name() string {
	return a.name
}

var (
	midiState *midi.MIDIState
	frame     PixelStateFrame
	ColorOff  = colorful.Color{R: 0, G: 0, B: 0}

	DefaultCommonSettings = &CommonSettings{
		LowerColor: colorful.FastLinearRgb(0, 1, 0),
		UpperColor: colorful.FastLinearRgb(1, 0, 0),
	}

	Animations = []control.ProcessInterface[midi.MIDIEvent, PixelStateFrame]{
		VelocityBar,
		VelocityBarMirror,
		FlowingNotes,
	}
)

func (ps *PixelState) Clear() {
	ps.Color = colorful.Color{R: 0, G: 0, B: 0}
	ps.Intensity = 0
}

func Initialize(rows int, cols int) {
	frame = PixelStateFrame{
		RowCount: rows,
		ColCount: cols,
		Pixels:   make([][]PixelState, rows),
	}
	for i := range frame.Pixels {
		frame.Pixels[i] = make([]PixelState, cols)
	}

	midiState = midi.NewMIDIState(cols, midi.PianoKeyboardDefault.KeyOffset)
}

func updateKeys(me midi.MIDIEvent) {
	k := midiState.Keys[me.KeyIndex]
	k.IsNotePressed = me.KeyPressed
	k.Velocity = me.Velocity
	if k.IsNotePressed {
		k.StartTime = me.TimeStamp
	} else {
		k.ReleaseTime = me.TimeStamp
	}
	midiState.UpdateVelocityRange(me.Velocity)
}
