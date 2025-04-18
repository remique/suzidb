package storage

import (
	"encoding/json"

	"example.com/suzidb/meta"
)

type Catalog interface {
	CreateTable(name string, schema meta.Table) error
	GetTable(name string) (*meta.Table, error)
	// TODO: Add []meta.Table as a key
	// ListTables() ([]string, error)
}

type SchemaManager struct {
	Storage Storage
}

func NewSchemaManager(storage Storage) *SchemaManager {
	return &SchemaManager{Storage: storage}
}

func (sm *SchemaManager) CreateTable(name string, schema meta.Table) error {
	serializedSchema, err := json.Marshal(schema)
	if err != nil {
		return err
	}
	sm.Storage.Set("meta:"+name, string(serializedSchema))

	return nil
}

func (sm *SchemaManager) GetTable(name string) (*meta.Table, error) {
	res := sm.Storage.Get("meta:" + name)
	if len(res) == 0 {
		return nil, nil
	}

	var deserializedSchema meta.Table
	err := json.Unmarshal([]byte(res), &deserializedSchema)
	if err != nil {
		return nil, err
	}

	return &deserializedSchema, nil
}
