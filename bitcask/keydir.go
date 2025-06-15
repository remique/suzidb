package bitcask

type KeyDirRecord struct {
	FileId    int
	ValueSize int
	ValuePos  int
	Timestamp int
}

// TODO: Use btreeMap?
type KeyDir map[string]KeyDirRecord

// TODO: We could implement methods on KeyDir as to simplify changing the underlying
// storage of KeyDir.
