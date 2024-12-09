package storage

type MemStorage struct {
	Store map[string]string
}

func NewMemStorage() *MemStorage {
	m := make(map[string]string)

	return &MemStorage{
		Store: m,
	}
}

func (ms *MemStorage) Get(key string) string {
	res, exist := ms.Store[key]
	if !exist {
		return ""
	}

	return res
}

func (ms *MemStorage) Set(key, value string) error {
	ms.Store[key] = value

	return nil
}
