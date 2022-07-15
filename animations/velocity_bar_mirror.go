package animations

import (
	"context"
	"math"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/whgentry/gomidi-led/leds"
)

var VelocityBarMirror = &Animation{
	Name:        "Velocity Bars Mirrored",
	Key:         "velocity-bars-mirror",
	Description: "Bars go up and down corresponding to the velocity of the note starting in the middle",
	Run: func(ctx context.Context) {
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
						} else if kboard.Keys[col].IsNotePressed {
							ps.Color = colorful.Hsv(360*float64(row)/float64(numRows), 1, 1)
						}
					}
				}
			case <-ctx.Done():
				return
			}
		}
	},
}
