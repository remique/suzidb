package storage

type KeyValue struct {
	Key   string
	Value string
}

type Storage interface {
	Get(key string) string
	Set(key, value string) error
}

// TODO: Add tuple saving
// TODO: Add catalog and table management
