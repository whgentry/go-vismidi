package animations

import (
	"context"

	"github.com/whgentry/gomidi-led/leds"
)

type VelocityBarAnimation struct {
}

func getLEDColor(row int, ps *PixelState) leds.Color {
	turnOnLED := row < int(ps.Intensity*float32(numRows))
	if turnOnLED {
		if row > 2*numRows/3 {
			return leds.Red
		} else if row > numRows/3 {
			return leds.Green
		} else {
			return leds.Blue
		}
	} else {
		return leds.Off
	}
}

func (a *VelocityBarAnimation) FrameHandler(lg leds.LEDGridInterface) {
	for i, ki := range kboard.Keys {
		for j := 0; j < numRows; j++ {
			lg.SetLED(j, i, getLEDColor(j, pixels[j][i]))
			if ki.IsNotePressed {
				pixels[j][i].Intensity = ki.GetAdjustedVelocityRatio()
			} else {
				pixels[j][i].Intensity *= 0.95
			}
		}
	}
}

func (a *VelocityBarAnimation) Run(ctx context.Context) {
}
