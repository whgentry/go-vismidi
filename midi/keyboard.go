package midi

import (
	"time"

	"gitlab.com/gomidi/midi/v2"
)

var PianoKeyboardDefault = &MIDIListener{
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

func NewMIDIState(keyCount int, keyOffset int) *MIDIState {
	ms := &MIDIState{
		Keys:        make([]*KeyInfo, keyCount),
		MaxVelocity: 51,
		MinVelocity: 50,
	}
	for i := range ms.Keys {
		ms.Keys[i] = &KeyInfo{
			MIDIState: ms,
			NoteName:  string(midi.Note(i + keyOffset)),
			StartTime: time.Now(),
			Index:     i,
		}
	}
	return ms
}

func (k *MIDIState) UpdateVelocityRange(vel int) {
	if vel > k.MaxVelocity {
		k.MaxVelocity = vel
	}
	if vel < k.MinVelocity {
		k.MinVelocity = vel
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
