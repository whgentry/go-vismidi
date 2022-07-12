package main

import (
	"context"
	"os"

	"github.com/nsf/termbox-go"
	"github.com/whgentry/gomidi-led/keyboard"
	"github.com/whgentry/gomidi-led/leds"
)

const NumLEDPerCol = 50
const midiPort = 0

var ledGrid leds.LEDGridInterface
var kboard *keyboard.Keyboard
var pixels [][]*PixelState

type PixelState struct {
	Color     leds.Color
	Intensity float32
}

func GetLEDColor(row int, ps *PixelState) leds.Color {
	turnOnLED := row < int(ps.Intensity*NumLEDPerCol)
	if turnOnLED {
		if row > 2*NumLEDPerCol/3 {
			return leds.Red
		} else if row > NumLEDPerCol/3 {
			return leds.Green
		} else {
			return leds.Blue
		}
	} else {
		return leds.Off
	}
}

func UpdateFrame(lg leds.LEDGridInterface) {
	for i, ki := range kboard.Keys {
		for j := 0; j < NumLEDPerCol; j++ {
			lg.SetLED(j, i, GetLEDColor(j, pixels[j][i]))
			if ki.IsNotePressed {
				pixels[j][i].Intensity = ki.GetAdjustedVelocityRatio()
			} else {
				pixels[j][i].Intensity *= 0.95
			}
		}
	}
}

// func LEDSmoothing(ctx context.Context) {
// 	ticker := time.NewTicker(10 * time.Millisecond)
// 	for {
// 		select {
// 		case <-ticker.C:
// 			for _, ki := range kboard.Keys {

// 			}
// 		case <-ctx.Done():
// 			return
// 		}
// 	}
// }

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	ledGrid = leds.NewVirtualLEDGrid(NumLEDPerCol, keyboard.KeyCount)
	kboard = keyboard.NewKeyboard()
	pixels = make([][]*PixelState, NumLEDPerCol)
	for i := range pixels {
		pixels[i] = make([]*PixelState, keyboard.KeyCount)
		for j := range pixels[i] {
			pixels[i][j] = &PixelState{}
		}
	}

	// If using virtual LED
	go leds.AnimateVirtualGrid(ctx, ledGrid.(*leds.VirtualLEDGrid), 60)

	// Core functionality routines
	go keyboard.HandleMidi(ctx, kboard, midiPort)
	go leds.HandleRefresh(ctx, ledGrid, 60, UpdateFrame)

	// Extra threads
	// go LEDSmoothing(ctx)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
				cancel()
				os.Exit(0)
			}
		}
	}
}
