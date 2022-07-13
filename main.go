package main

import (
	"context"
	"os"

	"github.com/nsf/termbox-go"
	"github.com/whgentry/gomidi-led/animations"
	"github.com/whgentry/gomidi-led/keyboard"
	"github.com/whgentry/gomidi-led/leds"
)

const NumLEDPerCol = 80
const midiPort = 0

var ledGrid leds.LEDGridInterface
var kboard *keyboard.Keyboard
var animationName string = "velocity-bar"
var frameRate = 60

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	if len(os.Args) > 1 {
		animationName = os.Args[1]
	}

	// Input and output structures
	ledGrid = leds.NewVirtualLEDGrid(NumLEDPerCol, keyboard.KeyCount)
	kboard = keyboard.NewKeyboard()

	// Animation handling
	animations.Initialize(NumLEDPerCol, keyboard.KeyCount, frameRate, kboard)
	animations.SetAnimation(animationName)

	// If using virtual LED
	go leds.AnimateVirtualGrid(ctx, ledGrid.(*leds.VirtualLEDGrid), frameRate)

	// Core functionality routines
	go keyboard.HandleMidi(ctx, kboard, midiPort)
	go leds.HandleRefresh(ctx, ledGrid, 60, animations.FrameHandler)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
				cancel()
				animations.Stop()
				os.Exit(0)
			}
		}
	}
}
