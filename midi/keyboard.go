package midi

import (
	"time"

	"gitlab.com/gomidi/midi/v2"
)

var PianoKeyboardDefault = MIDIListener{
	Port:      0,
	KeyOffset: 21,
	KeyCount:  88,
}

type KeyInfo struct {
	MIDIState     *MIDIState
	NoteName      string
	Index         int
	Velocity      int
	StartTime     time.Time
	ReleaseTime   time.Time
	IsNotePressed bool
}

type MIDIState struct {
	Keys        []*KeyInfo
	MaxVelocity int
	MinVelocity int
}

func NewMIDIState(ml MIDIListener) *MIDIState {
	ms := &MIDIState{
		Keys:        make([]*KeyInfo, ml.KeyCount),
		MaxVelocity: 51,
		MinVelocity: 50,
	}
	for i := range ms.Keys {
		ms.Keys[i] = &KeyInfo{
			MIDIState: ms,
			NoteName:  string(midi.Note(i + ml.KeyOffset)),
			StartTime: time.Now(),
			Index:     i,
		}
	}
	return ms
}

func (k *MIDIState) UpdateVelocityRange(vel uint8) {
	if vel > uint8(k.MaxVelocity) {
		k.MaxVelocity = int(vel)
	}
	if vel < uint8(k.MinVelocity) {
		k.MinVelocity = int(vel)
	}
}

func (k *MIDIState) GetVelocityRange() int {
	return k.MaxVelocity - k.MinVelocity
}

func (ki *KeyInfo) GetAdjustedVelocity() int {
	return ki.Velocity - ki.MIDIState.MinVelocity
}

func (ki *KeyInfo) GetAdjustedVelocityRatio() float64 {
	velocityRange := float64(ki.MIDIState.GetVelocityRange())
	adjustedVelocity := float64(ki.GetAdjustedVelocity())
	return adjustedVelocity / velocityRange
}
