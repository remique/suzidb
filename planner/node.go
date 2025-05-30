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

// Scan performs a basic scan of a Table.
type NodeScan struct {
	Table meta.Table
}

// Projection filters out columns that need to be queried.
type NodeProjection struct {
	Source      NodeQuery
	Expressions *[]parser.Expression
}

// TODO: Add support for building an actual plan.
type NestedLoopJoin struct {
	Left      NodeQuery
	Right     NodeQuery
	Predicate *parser.Expression
}

// NodeQuery marker trait implementations
func (ns *NodeScan) NodeQuery()        {}
func (np *NodeProjection) NodeQuery()  {}
func (nlj *NestedLoopJoin) NodeQuery() {}

type NodeBuilder struct {
	Catalog storage.Catalog
}

func NewNodeBuilder(c storage.Catalog) *NodeBuilder {
	return &NodeBuilder{Catalog: c}
}

func (nb *NodeBuilder) BuildNode(from parser.From) (NodeQuery, error) {
	switch f := from.(type) {
	case *parser.TableFrom:
		{
			return nb.buildNodeScan(f)
		}
	case *parser.JoinFrom:
		{
			left, err := nb.BuildNode(f.Left)
			if err != nil {
				return nil, err
			}

			right, err := nb.BuildNode(f.Right)
			if err != nil {
				return nil, err
			}
			return &NestedLoopJoin{
				Left:      left,
				Right:     right,
				Predicate: f.Predicate,
			}, nil
		}
	default:
		return nil, fmt.Errorf("Unsupported query")
	}
}

func (nb *NodeBuilder) buildNodeScan(tableFrom *parser.TableFrom) (NodeQuery, error) {
	// Get table
	table, err := nb.Catalog.GetTable(tableFrom.TableName)
	if err != nil {
		return nil, err
	}

	return &NodeScan{Table: *table}, nil
}

func (nb *NodeBuilder) AddNodeProjection(source NodeQuery, statement parser.Statement) (NodeQuery, error) {
	return &NodeProjection{
		Source:      source,
		Expressions: statement.SelectStatement.SelectItems,
	}, nil
}

func isAsteriskOnly(selectItems *[]lexer.Token) bool {
	if len(*selectItems) == 1 && (*selectItems)[0].TokenType == lexer.STAR {
		return true
	}

	return false
}
