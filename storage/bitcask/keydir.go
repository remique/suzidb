package bitcask

type KeyDirRecord struct {
	FileId    int
	ValueSize int
	ValuePos  int
	Timestamp int
}

// TODO: Use btreeMap?
type KeyDir map[string]KeyDirRecord
