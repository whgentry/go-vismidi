package animation

import (
	"context"
	"time"

	"github.com/whgentry/gomidi-led/midi"
)

var VelocityBar = &Animation{
	name:        "Velocity Bar",
	Description: "Bars go up corresponding to the velocity of the note",
	Settings:    velocityBarSettings,
	run:         velocityBarRun,
}

var velocityBarSettings = &Settings{
	CommonSettings: DefaultCommonSettings,
}

func velocityBarRun(ctx context.Context, input chan midi.MIDIEvent, out chan PixelStateFrame) {
	settings := velocityBarSettings
	frameTicker := time.NewTicker(50 * time.Millisecond)
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
					if row >= int(ps.Intensity*float64(frame.RowCount)) {
						ps.Color = ColorOff
					} else if midiState.Keys[col].IsNotePressed {
						ps.Color = settings.LowerColor.BlendHsv(settings.UpperColor, float64(row)/float64(frame.RowCount))
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
