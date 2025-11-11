package storage

type KeyValue struct {
	Key   string
	Value string
}

type Storage interface {
	Get(key string) string
	Set(key, value string) error
	ScanKeys() []string
	ScanWithPrefix(prefix string) map[string]string
}
