package keyboard

func MidiToKeyboardIndex(key uint8) int {
	if int(key) < MidiKeyboardOffset {
		return -1
	}
	return int(key) - MidiKeyboardOffset
}
