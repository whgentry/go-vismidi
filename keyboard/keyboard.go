package keyboard

import (
	"time"

	"gitlab.com/gomidi/midi/v2"
)

const MidiKeyboardOffset = 21
const KeyCount = 88

type KeyInfo struct {
	Keyboard      *Keyboard
	NoteName      string
	Index         int
	Velocity      int
	StartTime     time.Time
	ReleaseTime   time.Time
	IsNotePressed bool
}

type Keyboard struct {
	Keys        []*KeyInfo
	MaxVelocity int
	MinVelocity int
}

func NewKeyboard() *Keyboard {
	keyboard := &Keyboard{
		Keys:        make([]*KeyInfo, KeyCount),
		MaxVelocity: 51,
		MinVelocity: 50,
	}
	for i := range keyboard.Keys {
		keyboard.Keys[i] = &KeyInfo{
			Keyboard:  keyboard,
			NoteName:  string(midi.Note(i + MidiKeyboardOffset)),
			StartTime: time.Now(),
			Index:     i,
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

func (k *Keyboard) GetPressedKeys() []*KeyInfo {
	keys := []*KeyInfo{}
	for _, ki := range k.Keys {
		if ki.IsNotePressed {
			keys = append(keys, ki)
		}
	}
	return keys
}

func (k *Keyboard) GetVelocityRange() int {
	return k.MaxVelocity - k.MinVelocity
}

func (ki *KeyInfo) GetAdjustedVelocity() int {
	return ki.Velocity - ki.Keyboard.MinVelocity
}

func (ki *KeyInfo) GetAdjustedVelocityRatio() float32 {
	velocityRange := float32(ki.Keyboard.GetVelocityRange())
	adjustedVelocity := float32(ki.GetAdjustedVelocity())
	return adjustedVelocity / velocityRange
}
