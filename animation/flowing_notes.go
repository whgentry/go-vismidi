package animation

import (
	"context"
	"time"

	"github.com/whgentry/gomidi-led/midi"
)

var FlowingNotes = &Animation{
	name:        "Flowing Notes",
	Description: "Notes flow from the bottom upwards",
	Settings:    flowingNotesSettings,
	run:         flowingNotesRun,
}

var flowingNotesSettings = &Settings{
	CommonSettings: DefaultCommonSettings,
}

func flowingNotesRun(ctx context.Context, input chan midi.MIDIEvent, out chan PixelStateFrame) {
	settings := flowingNotesSettings
	speed := time.NewTicker(50 * time.Millisecond)
	for {
		select {
		case <-speed.C:
			for i := len(frame.Pixels) - 1; i >= 0; i-- {
				for j, ps := range frame.Pixels[i] {
					// Moves the bar up the column
					if i > 0 {
						ps.Intensity = frame.Pixels[i-1][j].Intensity
					} else {
						if midiState.Keys[j].IsNotePressed {
							ps.Intensity = midiState.Keys[j].GetAdjustedVelocityRatio()
						} else {
							ps.Intensity = 0
						}
					}
					// Sets pixel color based on intensity
					if ps.Intensity > 0 {
						ps.Color = settings.LowerColor.BlendHsv(settings.UpperColor, ps.Intensity)
					} else {
						ps.Color = ColorOff
					}
					frame.Pixels[i][j] = ps
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
