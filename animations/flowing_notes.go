package animations

import (
	"context"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/whgentry/gomidi-led/leds"
)

var FlowingNotes = &Animation{
	Name:        "Flowing Notes",
	Key:         "flowing-notes",
	Description: "Notes flow from the bottom upwards",
	Run: func(ctx context.Context) {
		defer wg.Done()
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
							ps.Color = colorful.Hsv(360*ps.Intensity, 1, 1)
						} else {
							ps.Color = leds.ColorOff()
						}
					}
				}
			case <-ctx.Done():
				return
			}
		}
	},
}
