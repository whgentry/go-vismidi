package animations

import (
	"context"
	"time"

	"github.com/whgentry/gomidi-led/leds"
)

var VelocityBar = &Animation{
	Name:        "Velocity Bars",
	Key:         "velocity-bars",
	Description: "Bars go up corresponding to the velocity of the note",
	Run: func(ctx context.Context, settings Settings) {
		defer wg.Done()
		frameTicker := time.NewTicker(frameDuration)
		for {
			select {
			case <-frameTicker.C:
				for row := range pixels {
					for col, ps := range pixels[row] {
						// Decay Intensity Exponentially
						if kboard.Keys[col].IsNotePressed {
							ps.Intensity = kboard.Keys[col].GetAdjustedVelocityRatio()
						} else {
							ps.Intensity *= 0.95
						}
						// Determine led color on intensity
						if row >= int(ps.Intensity*float64(numRows)) {
							ps.Color = leds.ColorOff()
						} else if kboard.Keys[col].IsNotePressed {
							ps.Color = settings.LowerColor.BlendHsv(settings.UpperColor, float64(row)/float64(numRows))
						}
					}
				}
			case <-ctx.Done():
				return
			}
		}
	},
}
