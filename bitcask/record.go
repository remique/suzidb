package bitcask

import (
	"encoding/json"
	"fmt"
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

func decode(input []byte) (*DiskRecord, error) {
	var dr DiskRecord
	err := json.Unmarshal(input, &dr)
	if err != nil {
		return nil, fmt.Errorf("Could not Unmarshal")
	}

	return &dr, nil
}
