package animations

import (
	"context"

	"github.com/whgentry/gomidi-led/leds"
)

type FlowingNotesAnimation struct {
}

func (a *FlowingNotesAnimation) FrameHandler(lg leds.LEDGridInterface) {
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

func (a *FlowingNotesAnimation) Run(ctx context.Context) {
}
