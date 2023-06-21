package midi

import (
	"context"
)

type MIDIReader struct {
	FileName string
}

var FileReaderDefault = &MIDIReader{}

func (m *MIDIReader) Name() string {
	return "MIDI Reader"
}

func (m *MIDIReader) Run(ctx context.Context, _ chan any, output chan MIDIEvent) {

}
