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
	animations.Initialize(NumLEDPerCol, midi.PianoKeyboardDefault.KeyCount)
	leds.Initialize(NumLEDPerCol, midi.PianoKeyboardDefault.KeyCount, frameRate)

	// Create Control Channels
	midiEventChan := make(chan midi.MIDIEvent)
	animationFrameChan := make(chan animations.PixelStateFrame)

	midiListener := midi.PianoKeyboardDefault
	midiCB := control.NewIOBlock(
		nil,
		midiEventChan,
		[]control.ProcessInterface[any, midi.MIDIEvent]{
			midiListener,
		},
	)

	animationCB := control.NewIOBlock(
		midiEventChan,
		animationFrameChan,
		animations.Animations,
	)

	ledCB := control.NewIOBlock(
		animationFrameChan,
		nil,
		leds.Displays,
	)

	// Start control blocks
	midiCB.Start(ctx)
	animationCB.Start(ctx)
	ledCB.Start(ctx)

	termbox.SetInputMode(termbox.InputEsc)
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowRight, termbox.KeySpace:
				animationCB.SetActive("Flowing Notes")
			case termbox.KeyArrowLeft:
				animationCB.SetActive("Velocity Bar")
			case termbox.KeyEsc, termbox.KeyCtrlC:
				os.Exit(0)
			}
		}
	}
}
