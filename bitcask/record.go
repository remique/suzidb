package bitcask

type DiskRecord struct {
	Header Header
	Key    string
	Value  []byte
}
