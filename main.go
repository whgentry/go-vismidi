package main

import (
	"context"
	"os"

	"github.com/nsf/termbox-go"
	"github.com/whgentry/gomidi-led/animations"
	"github.com/whgentry/gomidi-led/keyboard"
	"github.com/whgentry/gomidi-led/leds"
)

const NumLEDPerCol = 50
const midiPort = 0

var ledGrid leds.LEDGridInterface
var kboard *keyboard.Keyboard
var animation animations.Animation

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	animationctx, cancelAnimation := context.WithCancel(context.Background())

	// Input and output structures
	ledGrid = leds.NewVirtualLEDGrid(NumLEDPerCol, keyboard.KeyCount)
	kboard = keyboard.NewKeyboard()

	// Animation handling
	animations.Initialize(NumLEDPerCol, keyboard.KeyCount, kboard)
	animation = animations.GetAnimation("velocity-bar")
	animation.Run(animationctx)

	// If using virtual LED
	go leds.AnimateVirtualGrid(ctx, ledGrid.(*leds.VirtualLEDGrid), 60)

	// Core functionality routines
	go keyboard.HandleMidi(ctx, kboard, midiPort)
	go leds.HandleRefresh(ctx, ledGrid, 60, animation.FrameHandler)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
				cancel()
				cancelAnimation()
				os.Exit(0)
			}
		}
	}
}
