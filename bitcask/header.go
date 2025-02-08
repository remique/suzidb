package bitcask

import (
	"bytes"
	"hash/crc32"
	"time"
)

// Represents a header for a key-value
type Header struct {
	Crc       uint32
	Timestamp uint32
	KeySize   uint32
	ValueSize uint32
}

func NewHeader(key string, value []byte) *Header {
	return &Header{
		Crc:       crc32.ChecksumIEEE(bytes.NewBufferString(key).Bytes()),
		Timestamp: uint32(time.Now().Unix()),
		KeySize:   uint32(len(key)),
		ValueSize: uint32(len(value)),
	}
}
