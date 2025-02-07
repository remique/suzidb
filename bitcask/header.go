package bitcask

// Represents a header for a key-value
type Header struct {
	Crc       uint32
	Timestamp uint32
	KeySize   uint32
	ValueSize uint32
}

func (h *Header) encode() {
	//
}

func (h *Header) decode() {
	//
}
