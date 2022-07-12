package keyboard

import "gitlab.com/gomidi/midi/v2"

const MidiKeyboardOffset = 21
const KeyboardKeyCount = 88

type KeyInfo struct {
	NoteName      string
	Velocity      int
	IsNotePressed bool
}

type Keyboard struct {
	Keys        []*KeyInfo
	MaxVelocity int
	MinVelocity int
}

func NewKeyboard() *Keyboard {
	keyboard := &Keyboard{
		Keys:        make([]*KeyInfo, KeyboardKeyCount),
		MaxVelocity: 51,
		MinVelocity: 50,
	}
	for i := range keyboard.Keys {
		keyboard.Keys[i] = &KeyInfo{
			NoteName: string(midi.Note(i + MidiKeyboardOffset)),
		}
	}
	return keyboard
}

func (k *Keyboard) UpdateVelocityRange(vel uint8) {
	if vel > uint8(k.MaxVelocity) {
		k.MaxVelocity = int(vel)
	}
	if vel < uint8(k.MinVelocity) {
		k.MinVelocity = int(vel)
	}
}
