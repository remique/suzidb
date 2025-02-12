package bitcask

import (
	"encoding/json"
)

type DiskRecord struct {
	Header Header
	Key    string
	Value  []byte
}

func NewDiskRecord(key string, value string) *DiskRecord {
	return &DiskRecord{
		Header: *NewHeader(key, []byte(value)),
		Key:    key,
		Value:  []byte(value),
	}
}

func (dr *DiskRecord) encode() ([]byte, error) {
	serialized, err := json.Marshal(dr)
	if err != nil {
		return []byte{}, err
	}

	return serialized, nil
}
