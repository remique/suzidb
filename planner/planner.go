package planner

import (
	"fmt"
	"reflect"

	l "example.com/suzidb/lexer"
	m "example.com/suzidb/meta"
	p "example.com/suzidb/parser"
	s "example.com/suzidb/storage"
)

type Planner struct {
	Catalog s.Catalog
}

func (pl *Planner) Build(statement p.Statement) (Plan, error) {
	switch statement.Kind {
	case p.CreateTableKind:
		return pl.buildCreateTable(statement)
	case p.InsertKind:
		return pl.buildInsert(statement)
	case p.SelectKind:
		return pl.buildSelect(statement)
	}

	return nil, nil
}

func NewPlanner(c s.Catalog) *Planner {
	return &Planner{Catalog: c}
}

func (pl *Planner) buildCreateTable(statement p.Statement) (Plan, error) {
	tableExists, err := pl.Catalog.GetTable(statement.CreateTableStatement.TableName)
	if err != nil {
		return nil, fmt.Errorf("Error while fetching getTable: %s", err.Error())
	}
	if tableExists != nil {
		return nil, fmt.Errorf("Table already exists")
	}

	table := m.Table{
		Name:       statement.CreateTableStatement.TableName,
		Columns:    *statement.CreateTableStatement.Columns,
		PrimaryKey: statement.CreateTableStatement.PrimaryKey,
	}

	plan := CreateTablePlan{Table: table}

	return &plan, nil
}

func (pl *Planner) buildInsert(stmt p.Statement) (Plan, error) {
	row := make(map[string]interface{})

	// Get table
	tableName := stmt.InsertStatement.TableName
	table, err := pl.Catalog.GetTable(tableName)
	if err != nil {
		return nil, fmt.Errorf("Error while fetching getTable: %s", err.Error())
	}
	if table == nil {
		return nil, fmt.Errorf("Table %s does not exist", tableName)
	}

	if len(stmt.InsertStatement.CustomColumns) != len(*stmt.InsertStatement.Values) &&
		len(stmt.InsertStatement.CustomColumns) > 0 {
		return nil, fmt.Errorf("Got %d columns and %d values",
			len(stmt.InsertStatement.CustomColumns), len(*stmt.InsertStatement.Values))
	}

	// TODO: Refactor this
	if len(stmt.InsertStatement.CustomColumns) == 0 {
		for i, c := range table.Columns {
			// Check the type
			currExpr := (*stmt.InsertStatement.Values)[i]
			if c.Type != exprToColumnType(&currExpr) {
				return nil, fmt.Errorf("Expected %d, got %d", c.Type, exprToColumnType(&currExpr))
			}

			toToken, err := exprToToken(&currExpr)
			if err != nil {
				return nil, fmt.Errorf("Error toToken")
			}

			row[c.Name] = toToken.Literal
		}
	} else {
		for _, c := range table.Columns {
			// Get index of customCols
			idx := getColumnIndex(stmt.InsertStatement.CustomColumns, c.Name)
			if idx == -1 {
				if c.Nullable == true {
					row[c.Name] = ""
				} else {
					return nil, fmt.Errorf("Error while getting column")
				}
			} else {
				currExpr := (*stmt.InsertStatement.Values)[idx]
				if c.Type != exprToColumnType(&currExpr) {
					return nil, fmt.Errorf("Expected %d, got %d", c.Type, exprToColumnType(&currExpr))
				}

				toToken, err := exprToToken(&currExpr)
				if err != nil {
					return nil, fmt.Errorf("Error toToken")
				}

				row[c.Name] = toToken.Literal
			}
		}
	}

	return &InsertPlan{Table: *table, Row: row}, nil
}

func (pl *Planner) buildSelect(stmt p.Statement) (Plan, error) {
	nb := NewNodeBuilder(pl.Catalog)

	node, err := nb.BuildNode(stmt.SelectStatement.From)
	if err != nil {
		return nil, err
	}

	// Here we can build projection if it exists
	if !reflect.DeepEqual(stmt.SelectStatement.SelectItems, &[]p.Expression{{Kind: p.AllColumnsKind, AllColumnsExpression: &p.AllExpression{}}}) {
		node, err = nb.AddNodeProjection(node, stmt)
		if err != nil {
			return nil, err
		}
	}

	return &SelectPlan{Node: node}, nil
}

func getColumnIndex(slice []l.Token, columnName string) int {
	for i := range slice {
		if columnName == slice[i].Literal {
			return i
		}
	}

	return -1
}

// func tokenToColumnType(token l.Token) m.ColumnType {
// 	// NOTE(remique): We need to clean the types and define them
// 	// more thoroughly.
// 	switch token.TokenType {
// 	case l.STRING:
// 		return m.StringType
// 	case l.INT:
// 		return m.IntType
// 	case l.TEXT_TYPE:
// 		return m.StringType
// 	case l.INT_TYPE:
// 		return m.IntType
// 	default:
// 		return -1
// 	}
// }

func exprToColumnType(expr *p.Expression) m.ColumnType {

	switch expr.Kind {
	case p.LiteralKind:
		return m.StringType
	case p.ConstExprKind:
		{
			switch expr.ConstExpression.Kind {
			case p.IntKind:
				return m.IntType
			case p.StringKind:
				return m.StringType
			default:
				return -1
			}
		}
	default:
		return -1
	}
}

func exprToToken(expr *p.Expression) (*l.Token, error) {
	switch expr.Kind {
	case p.LiteralKind:
		{
			return expr.LiteralExpression, nil
		}
	case p.ConstExprKind:
		{
			switch expr.ConstExpression.Kind {
			case p.IntKind:
				return expr.ConstExpression.Int, nil
			default:
				return nil, fmt.Errorf("Unsupported ConstExpr")
			}
		}
	default:
		{
			return nil, fmt.Errorf("Unsupported Expression")
		}
	}
}
