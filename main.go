package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/whgentry/gomidi-led/leds"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

type KeyInfo struct {
	Note     string
	Velocity int
}

const MidiKeyboardOffset = 21
const KeyboardKeyCount = 88
const NumLEDPerCol = 50

var ledGrid leds.LEDGridInterface
var keyboardKeys []KeyInfo

var MaxVelocity = 51
var MinVelocity = 50

func MidiToKeyboardIndex(key uint8) int {
	if int(key) < MidiKeyboardOffset {
		return -1
	}
	return int(key) - MidiKeyboardOffset
}

func UpdateVelocityRange(vel uint8) {
	if vel > uint8(MaxVelocity) {
		MaxVelocity = int(vel)
	}
	if vel < uint8(MinVelocity) {
		MinVelocity = int(vel)
	}
}

func GetLEDColor(row int, ki KeyInfo) leds.Color {
	velocityRange := MaxVelocity - MinVelocity
	adjustedVelocity := ki.Velocity - MinVelocity
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

func UpdateFrame(ctx context.Context, refreshRate int) {
	frameDurationMs := 1000 / refreshRate
	ticker := time.NewTicker(time.Duration(frameDurationMs) * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			for i, keyInfo := range keyboardKeys {
				for j := 0; j < NumLEDPerCol; j++ {
					ledGrid.SetLED(j, i, GetLEDColor(j, keyInfo))
				}
			}
			ledGrid.UpdateLEDs()
		case <-ctx.Done():
			return
		}
	}

}

func main() {
	defer midi.CloseDriver()

	ledGrid = leds.NewVirtualLEDGrid(NumLEDPerCol, KeyboardKeyCount)
	keyboardKeys = make([]KeyInfo, KeyboardKeyCount)

	ctx, cancel := context.WithCancel(context.Background())

	in, err := midi.InPort(0)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	fmt.Print("MIDI Device Connected")

	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var bt []byte
		var ch, key, vel uint8
		switch {
		case msg.GetSysEx(&bt):
			fmt.Printf("got sysex: % X\n", bt)
		case msg.GetNoteStart(&ch, &key, &vel):
			keyboardKeys[MidiToKeyboardIndex(key)].Velocity = int(vel)
			UpdateVelocityRange(vel)
		case msg.GetNoteEnd(&ch, &key):
			keyboardKeys[MidiToKeyboardIndex(key)].Velocity = 0
		default:
			// ignore
		}
	}, midi.UseSysEx())

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	go leds.AnimateVirtualGrid(ledGrid.(*leds.VirtualLEDGrid), 30)
	go UpdateFrame(ctx, 30)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
				cancel()
				leds.StopVirtualGrid(ledGrid.(*leds.VirtualLEDGrid))
				stop()
				os.Exit(0)
			}
		}
	}
}
