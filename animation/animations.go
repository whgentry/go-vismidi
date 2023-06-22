package animation

import (
	"context"

	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/whgentry/go-vismidi/control"
	"github.com/whgentry/go-vismidi/midi"
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

type Animation struct {
	Description string
	Settings    *AnimationSettings
	name        string
	run         func(context.Context, chan midi.MIDIEvent, chan PixelStateFrame, *AnimationSettings)
}

func (a *Animation) Run(ctx context.Context, input chan midi.MIDIEvent, out chan PixelStateFrame) {
	a.run(ctx, input, out, a.Settings)
}

func (a *Animation) Name() string {
	return a.name
}

var (
	midiState *midi.MIDIState
	frame     PixelStateFrame
	ColorOff  = colorful.Color{R: 0, G: 0, B: 0}

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
