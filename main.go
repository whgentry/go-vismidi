package main

import (
	"context"
	"os"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/whgentry/gomidi-led/keyboard"
	"github.com/whgentry/gomidi-led/leds"
)

const NumLEDPerCol = 50
const midiPort = 0

var ledGrid leds.LEDGridInterface
var kboard *keyboard.Keyboard

func GetLEDColor(row int, ki *keyboard.KeyInfo) leds.Color {
	velocityRange := kboard.MaxVelocity - kboard.MinVelocity
	adjustedVelocity := ki.Velocity - kboard.MinVelocity
	turnOnLED := row < (adjustedVelocity*NumLEDPerCol)/velocityRange
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
	for i, keyInfo := range kboard.Keys {
		for j := 0; j < NumLEDPerCol; j++ {
			lg.SetLED(j, i, GetLEDColor(j, keyInfo))
		}
	}
}

func LEDSmoothing(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			for _, ki := range kboard.Keys {
				if !ki.IsNotePressed {
					ki.Velocity = int(float32(ki.Velocity) * 0.98)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	ledGrid = leds.NewVirtualLEDGrid(NumLEDPerCol, keyboard.KeyboardKeyCount)
	kboard = keyboard.NewKeyboard()

	// If using virtual LED
	go leds.AnimateVirtualGrid(ctx, ledGrid.(*leds.VirtualLEDGrid), 60)

	go keyboard.HandleMidi(ctx, kboard, midiPort)
	go leds.HandleRefresh(ctx, ledGrid, 60, UpdateFrame)

	go LEDSmoothing(ctx)

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
