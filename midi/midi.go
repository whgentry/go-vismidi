package midi

import (
	"context"
	"fmt"
	"time"

	"github.com/whgentry/go-vismidi/control"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

type MIDIEvent struct {
	Velocity   int
	KeyIndex   int
	KeyPressed bool
	TimeStamp  time.Time
}

type MIDIListener struct {
	Settings
}

var Inputs = []control.ProcessInterface[any, MIDIEvent]{
	PianoKeyboardDefault,
	FileReaderDefault,
}

func (m *MIDIListener) MidiToKeyboardIndex(key uint8) int {
	if int(key) < m.KeyOffset {
		return -1
	}
	return int(key) - m.KeyOffset
}

func (m *MIDIListener) Name() string {
	return "MIDI Listener"
}

func (m *MIDIListener) Run(ctx context.Context, _ chan any, output chan MIDIEvent) {
	defer midi.CloseDriver()

	var in drivers.In
	var err error

	for in, err = midi.InPort(m.Port); err != nil; in, err = midi.InPort(m.Port) {
	}

	stopListening, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var bt []byte
		var ch, key, vel uint8
		switch {
		case msg.GetSysEx(&bt):
			fmt.Printf("got sysex: % X\n", bt)
		case msg.GetNoteStart(&ch, &key, &vel):
			output <- MIDIEvent{
				Velocity:   int(vel),
				KeyIndex:   m.MidiToKeyboardIndex(key),
				TimeStamp:  time.Now(),
				KeyPressed: true,
			}
		case msg.GetNoteEnd(&ch, &key):
			output <- MIDIEvent{
				KeyIndex:   m.MidiToKeyboardIndex(key),
				TimeStamp:  time.Now(),
				KeyPressed: false,
			}
		default:
			// ignore
		}
	}, midi.UseSysEx())

	if err != nil {
		fmt.Printf("ERROR: Failed to start midi listening %v\n", err)
		return
	}

	for {
		<-ctx.Done()
		stopListening()
		return
	}
}
