package animation

import (
	"context"
	"math"
	"time"

	"github.com/whgentry/gomidi-led/midi"
)

var VelocityBarMirror = &Animation{
	name:        "Velocity Bar Mirrored",
	Description: "Bars go up and down corresponding to the velocity of the note starting in the middle",
	Settings:    velocityBarMirrorSettings,
	run:         velocityBarMirrorRun,
}

var velocityBarMirrorSettings = &Settings{
	CommonSettings: DefaultCommonSettings,
}

func velocityBarMirrorRun(ctx context.Context, input chan midi.MIDIEvent, out chan PixelStateFrame) {
	settings := velocityBarMirrorSettings
	frameTicker := time.NewTicker(50 * time.Millisecond)
	middleRow := frame.RowCount / 2
	for {
		select {
		case <-frameTicker.C:
			for row := range frame.Pixels {
				for col, ps := range frame.Pixels[row] {
					// Decay Intensity Exponentially
					if midiState.Keys[col].IsNotePressed {
						ps.Intensity = midiState.Keys[col].GetAdjustedVelocityRatio()
					} else {
						ps.Intensity *= 0.95
					}
					// Determine led color on intensity
					if int(math.Abs(float64(middleRow-row))) >= int(ps.Intensity*float64(middleRow)) {
						ps.Color = ColorOff
					} else {
						ps.Color = settings.LowerColor.BlendHsv(settings.UpperColor, ps.Intensity)
					}
					frame.Pixels[row][col] = ps
				}
			}
			out <- frame
		case me := <-input:
			updateKeys(me)
		case <-ctx.Done():
			return
		}
	}
}
