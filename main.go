package main

import (
	"context"
	"os"

	"github.com/nsf/termbox-go"
	"github.com/whgentry/gomidi-led/animations"
	"github.com/whgentry/gomidi-led/control"
	"github.com/whgentry/gomidi-led/leds"
	"github.com/whgentry/gomidi-led/midi"
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

	// Create Control Channels
	midiEventChan := make(chan midi.MIDIEvent)

	midiListener := midi.PianoKeyboardDefault

	midiCB := control.IOBlock[any, midi.MIDIEvent]{
		Input:  nil,
		Output: midiEventChan,
		Processors: []control.ProcessInterface[any, midi.MIDIEvent]{
			midiListener,
		},
	}

	animationCB := control.IOBlock[midi.MIDIEvent, any]{
		Input:      midiEventChan,
		Output:     nil,
		Processors: animations.Animations,
	}

	midiCB.Start(ctx)
	animationCB.Start(ctx)

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
