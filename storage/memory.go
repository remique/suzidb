package storage

import (
	"strings"
)

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

func (ms *MemStorage) ScanKeys() []string {
	var res []string

	for key := range ms.Store {
		res = append(res, key)
	}

	return res
}

func (ms *MemStorage) ScanWithPrefix(prefix string) map[string]string {
	m := make(map[string]string)

	for key, value := range ms.Store {
		if strings.HasPrefix(key, prefix) {
			m[key] = value
		}
	}

	return m
}
