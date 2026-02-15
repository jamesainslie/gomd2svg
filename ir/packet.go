package ir

// PacketField represents a single field in a network packet header diagram.
type PacketField struct {
	Start       int
	End         int
	Description string
}

// BitWidth returns the number of bits this field spans.
func (f *PacketField) BitWidth() int {
	return f.End - f.Start + 1
}
