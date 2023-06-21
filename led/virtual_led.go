package led

import (
	"context"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/whgentry/gomidi-led/animation"
)

type VirtualLEDGrid struct {
	state animation.PixelStateFrame
}

func (lg *VirtualLEDGrid) Run(ctx context.Context, input chan animation.PixelStateFrame, out chan any) {
	ticker := time.NewTicker(time.Second / time.Duration(frameRate))
	for {
		select {
		case state := <-input:
			lg.state = state
		case <-ticker.C:
			termbox.SetOutputMode(termbox.OutputRGB)
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			for i := range lg.state.Pixels {
				for j, led := range lg.state.Pixels[i] {
					if !IsColorOff(led.Color) {
						r, g, b := led.Color.RGB255()
						fg := termbox.RGBToAttribute(r, g, b)
						bg := termbox.Attribute(termbox.ColorDefault)
						_, height := termbox.Size()
						termbox.SetCell(j*2, height-i, '*', fg, bg)
						termbox.SetCell(j*2+1, height-i, '*', fg, bg)
					}
				}
			}
			termbox.Flush()
		case <-ctx.Done():
			return
		}
	}
}

func (lg *VirtualLEDGrid) Name() string {
	return "VirtualLEDGrid"
}
