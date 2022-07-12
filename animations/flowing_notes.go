package animations

import (
	"context"
	"time"

	"github.com/whgentry/gomidi-led/leds"
)

type FlowingNotesAnimation struct {
}

func flowingNotesGetLEDColor(ps *PixelState) leds.Color {
	if ps.Intensity > 0 {
		if ps.Intensity > 0.66 {
			return leds.Red
		} else if ps.Intensity > 0.33 {
			return leds.Green
		} else {
			return leds.Blue
		}
	} else {
		return leds.Off
	}
}

func (a *FlowingNotesAnimation) FrameHandler(lg leds.LEDGridInterface) {
	for i := range pixels {
		for j := range pixels[i] {
			lg.SetLED(i, j, flowingNotesGetLEDColor(pixels[i][j]))
		}
	}
}

func (a *FlowingNotesAnimation) Run(ctx context.Context) {
	go flowingNotesMoveBar(ctx)
	<-ctx.Done()
}

func flowingNotesMoveBar(ctx context.Context) {
	speed := time.NewTicker(50 * time.Millisecond)
	for {
		select {
		case <-speed.C:
			for i := len(pixels) - 1; i >= 0; i-- {
				for j := range pixels[i] {
					if i > 0 {
						pixels[i][j].Intensity = pixels[i-1][j].Intensity
					} else {
						if kboard.Keys[j].IsNotePressed {
							pixels[i][j].Intensity = kboard.Keys[j].GetAdjustedVelocityRatio()
						} else {
							pixels[i][j].Intensity = 0
						}
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
