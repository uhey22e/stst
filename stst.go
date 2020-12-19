package stst

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

type Stst struct {
	db *sql.DB
}

var (
	dsn         = "postgresql://postgres@localhost:15432/postgres?sslmode=disable"
	errCols     = "Failed to read columns: %w"
	errColTypes = "Failed to read column types: %w"
)

// New is a constructor
func New(db *sql.DB) *Stst {
	return &Stst{
		db: db,
	}
}

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
func (s *Stst) GenerateStruct(cols [][2]string) (*jen.Statement, error) {
	// Use []jen.Code instead of []*jen.Statement to pass it to jen.Structs()
	ms := make([]jen.Code, len(cols))
	for i, c := range cols {
		ms[i] = jen.Id(strcase.ToCamel(c[0])).Add(jen.Id(c[1]))
	}

	st := jen.Type().Id("Foo").Struct(ms...)
	return st, nil
}

// Package .
func (s *Stst) Package(w io.Writer, name string, codes []jen.Code) error {
	f := jen.NewFile(name)
	for _, c := range codes {
		f.Add(c)
	}
	if err := f.Render(w); err != nil {
		return err
	}
	return nil
}
