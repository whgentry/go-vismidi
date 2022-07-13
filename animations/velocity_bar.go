package animations

import (
	"context"
	"time"

	"github.com/whgentry/gomidi-led/leds"
)

type VelocityBarAnimation struct {
}

func (a *VelocityBarAnimation) Run(ctx context.Context) {
	frameTicker := time.NewTicker(frameDuration)
	for {
		select {
		case <-frameTicker.C:
			for row := range pixels {
				for col, ps := range pixels[row] {
					// Determine led color on intensity
					if row < int(ps.Intensity*float32(numRows)) {
						if row > 2*numRows/3 {
							ps.Color = leds.Red
						} else if row > numRows/3 {
							ps.Color = leds.Green
						} else {
							ps.Color = leds.Blue
						}
					} else {
						ps.Color = leds.Off
					}
					// Decay Intensity Exponentially
					if kboard.Keys[col].IsNotePressed {
						ps.Intensity = kboard.Keys[col].GetAdjustedVelocityRatio()
					} else {
						ps.Intensity *= 0.95
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
