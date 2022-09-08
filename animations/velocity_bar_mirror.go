package animations

import (
	"context"
	"math"
	"time"

	"github.com/whgentry/gomidi-led/leds"
)

var VelocityBarMirror = &Animation{
	Name:        "Velocity Bars Mirrored",
	Description: "Bars go up and down corresponding to the velocity of the note starting in the middle",
	Run: func(ctx context.Context, settings Settings) {
		defer wg.Done()
		frameTicker := time.NewTicker(frameDuration)
		middleRow := numRows / 2
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
						if int(math.Abs(float64(middleRow-row))) >= int(ps.Intensity*float64(middleRow)) {
							ps.Color = leds.ColorOff()
						} else {
							ps.Color = settings.LowerColor.BlendHsv(settings.UpperColor, ps.Intensity)
						}
					}
				}
			case <-ctx.Done():
				return
			}
		}
	},
}
