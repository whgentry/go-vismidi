package main

import (
	"context"
	"os"

	"github.com/nsf/termbox-go"
	"github.com/whgentry/gomidi-led/animations"
	"github.com/whgentry/gomidi-led/keyboard"
	"github.com/whgentry/gomidi-led/leds"
)

var midiPort = 0
var NumLEDPerCol = 70
var ledGrid leds.LEDGridInterface
var kboard *keyboard.Keyboard
var frameRate = 60

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := termbox.Init()
	if err != nil {
		os.Exit(1)
	}
	defer termbox.Close()

	_, NumLEDPerCol = termbox.Size()

	// Input and output structures
	ledGrid = leds.NewVirtualLEDGrid(NumLEDPerCol, keyboard.KeyCount)
	kboard = keyboard.NewKeyboard()

	// Animation handling
	animations.Initialize(NumLEDPerCol, keyboard.KeyCount, frameRate, kboard)
	defer animations.Close()
	animations.SetAnimationByIndex(0)

	// Core functionality routines
	go keyboard.HandleMidi(ctx, kboard, midiPort)
	go leds.HandleRefresh(ctx, ledGrid, 60, animations.FrameHandler)

	termbox.SetInputMode(termbox.InputEsc)
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowRight, termbox.KeySpace:
				animations.NextAnimation()
			case termbox.KeyArrowLeft:
				animations.PreviousAnimation()
			case termbox.KeyEsc, termbox.KeyCtrlC:
				os.Exit(0)
			}
		}
	}
}
