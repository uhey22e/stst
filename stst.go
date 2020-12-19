package stst

import (
	"database/sql"
	"fmt"

	"github.com/dave/jennifer/jen"
)

type Stst struct {
	db *sql.DB
}

var (
	dsn         = "postgresql://postgres@localhost:15432/postgres?sslmode=disable"
	errCols     = "Failed to read columns: %w"
	errColTypes = "Failed to read column types: %w"
)

// GetMeta returns metadata of columns
func (s *Stst) GetMeta(query string) ([]string, []*sql.ColumnType, error) {
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, fmt.Errorf(errCols, err)
	}

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, nil, fmt.Errorf(errColTypes, err)
	}

	return cols, colTypes, nil
}

// GenerateStruct .
func (s *Stst) GenerateStruct() (*jen.Statement, error) {
	st := jen.Type().Id("Foo").Struct(
		jen.Id("id").Int().Tag(map[string]string{"json": "jsonKey"}),
	)
	return st, nil
}
