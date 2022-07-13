package animations

import (
	"context"
	"time"

	"github.com/whgentry/gomidi-led/leds"
)

type FlowingNotesAnimation struct {
}

func (a *FlowingNotesAnimation) Run(ctx context.Context) {
	speed := time.NewTicker(50 * time.Millisecond)
	for {
		select {
		case <-speed.C:
			for i := len(pixels) - 1; i >= 0; i-- {
				for j, ps := range pixels[i] {
					// Moves the bar up the column
					if i > 0 {
						pixels[i][j].Intensity = pixels[i-1][j].Intensity
					} else {
						if kboard.Keys[j].IsNotePressed {
							pixels[i][j].Intensity = kboard.Keys[j].GetAdjustedVelocityRatio()
						} else {
							pixels[i][j].Intensity = 0
						}
					}
					// Sets pixel color based on intensity
					if ps.Intensity > 0 {
						if ps.Intensity > 0.66 {
							ps.Color = leds.Red
						} else if ps.Intensity > 0.33 {
							ps.Color = leds.Green
						} else {
							ps.Color = leds.Blue
						}
					} else {
						ps.Color = leds.Off
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
