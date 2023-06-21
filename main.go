package main

import (
	"context"
	"os"

	"github.com/nsf/termbox-go"
	"github.com/whgentry/go-vismidi/animation"
	"github.com/whgentry/go-vismidi/control"
	"github.com/whgentry/go-vismidi/led"
	"github.com/whgentry/go-vismidi/midi"
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
	animation.Initialize(NumLEDPerCol, midi.PianoKeyboardDefault.KeyCount)
	led.Initialize(NumLEDPerCol, midi.PianoKeyboardDefault.KeyCount, frameRate)

	// Create Control Channels
	midiEventChan := make(chan midi.MIDIEvent, 100)
	animationFrameChan := make(chan animation.PixelStateFrame, 100)

	midiCB := control.NewIOBlock(
		nil,
		midiEventChan,
		midi.Inputs,
	)

	animationCB := control.NewIOBlock(
		midiEventChan,
		animationFrameChan,
		animation.Animations,
	)

	ledCB := control.NewIOBlock(
		animationFrameChan,
		nil,
		led.Displays,
	)

	// Start control blocks
	midiCB.Start(ctx)
	animationCB.Start(ctx)
	ledCB.Start(ctx)

	// TODO Add methods to controlblock to allow rotation through processor
	termbox.SetInputMode(termbox.InputEsc)
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowRight, termbox.KeySpace:
				animationCB.Next()
			case termbox.KeyArrowLeft:
				animationCB.Previous()
			case termbox.KeyEsc, termbox.KeyCtrlC:
				os.Exit(0)
			}
		}
	}
}
