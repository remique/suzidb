package mocks

import (
	"example.com/suzidb/meta"
)

type MockCatalog struct {
	CreateTableFunc func(name string, schema meta.Table) error
	GetTableFunc    func(name string) (*meta.Table, error)
}

func (mc *MockCatalog) CreateTable(name string, schema meta.Table) error {
	return mc.CreateTableFunc(name, schema)
}

func (mc *MockCatalog) GetTable(name string) (*meta.Table, error) {
	return mc.GetTableFunc(name)
}
