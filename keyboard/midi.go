package keyboard

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func HandleMidi(ctx context.Context, kboard *Keyboard, port int) {
	defer midi.CloseDriver()

	var in drivers.In
	var err error

	for in, err = midi.InPort(port); err != nil; in, err = midi.InPort(port) {
	}

	stopListening, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var bt []byte
		var ch, key, vel uint8
		switch {
		case msg.GetSysEx(&bt):
			fmt.Printf("got sysex: % X\n", bt)
		case msg.GetNoteStart(&ch, &key, &vel):
			kboard.Keys[MidiToKeyboardIndex(key)].Velocity = int(vel)
			kboard.Keys[MidiToKeyboardIndex(key)].IsNotePressed = true
			kboard.Keys[MidiToKeyboardIndex(key)].StartTime = time.Now()
			kboard.UpdateVelocityRange(vel)
		case msg.GetNoteEnd(&ch, &key):
			kboard.Keys[MidiToKeyboardIndex(key)].IsNotePressed = false
			kboard.Keys[MidiToKeyboardIndex(key)].ReleaseTime = time.Now()
		default:
			// ignore
		}
	}, midi.UseSysEx())

	if err != nil {
		fmt.Printf("ERROR: Failed to start midi listening %v\n", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			stopListening()
			return
		}
	}
}
