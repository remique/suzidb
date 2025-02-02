package planner

import (
	"fmt"

	"example.com/suzidb/lexer"
	"example.com/suzidb/meta"
	"example.com/suzidb/parser"
	"example.com/suzidb/storage"
)

// NodeQuery is a part of a SelectPlan, which will later on be passed to
// and executor. As can be seen below, NodeQuery may be recursive, eg.
// NodeProjection can have a source of NodeScan. Later on it might have a
// source of IndexLookup.
type NodeQuery interface {
	NodeQuery()
}

// Scan performs a basic scan of a Table. Currently, filtering out
// specific rows is not supported (eg. 'WHERE x > 5'), since we do not have
// proper expression parsing.
type NodeScan struct {
	Table meta.Table

	// filter: parser.Expression
}

// Projection filters out columns that need to be queried.
type NodeProjection struct {
	Source  NodeQuery
	Columns []meta.Column
}

// NodeQuery marker trait implementations
func (ns *NodeScan) NodeQuery()       {}
func (np *NodeProjection) NodeQuery() {}

type NodeBuilder struct {
	Catalog storage.Catalog
}

func NewNodeBuilder(c storage.Catalog) *NodeBuilder {
	return &NodeBuilder{Catalog: c}
}

func (nb *NodeBuilder) BuildNode(statement parser.Statement) (NodeQuery, error) {
	switch statement.Kind {
	case parser.SelectKind:
		{
			if isAsteriskOnly(statement.SelectStatement.SelectItems) {
				return nb.buildNodeScan(statement)
			}

			return nil, fmt.Errorf("Unsupported query")
		}
	default:
		return nil, fmt.Errorf("Expected SelectKind")
	}
}

func (nb *NodeBuilder) buildNodeScan(statement parser.Statement) (NodeQuery, error) {
	// Get table
	table, err := nb.Catalog.GetTable(statement.SelectStatement.From.Literal)
	if err != nil {
		return nil, err
	}

	return &NodeScan{Table: *table}, nil
}

// func (nb *NodeBuilder) buildNodeProjection(statement parser.Statement) (NodeQuery, error) {
// 	// Build source
//      // source := buildNodeScan()
//      // columns: Get columns from statement

//      return &NodeProjection{...}

// 	return &NodeScan{Table: *table}, nil
// }

func isAsteriskOnly(selectItems *[]lexer.Token) bool {
	if len(*selectItems) == 1 && (*selectItems)[0].TokenType == lexer.STAR {
		return true
	}

	return false
}
