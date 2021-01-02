package stst

import (
	"database/sql"
	"fmt"
	"io"
	"reflect"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

// Stst .
type Stst struct {
	DB      *sql.DB
	Typemap Typemap
}

var (
	dsn         = "postgresql://postgres@localhost:15432/postgres?sslmode=disable"
	errCols     = "Failed to read columns"
	errColTypes = "Failed to read column types"
)

// NewPsql is a constructor
func NewPsql(db *sql.DB) *Stst {
	return &Stst{
		DB:      db,
		Typemap: PsqlTypemap{},
	}
}

// GetMeta returns metadata of columns
func (s *Stst) GetMeta(query string) ([]string, []string, error) {
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", errCols, err)
	}

	colTypes := make([]string, len(cols))
	cts, err := rows.ColumnTypes()
	for i, ct := range cts {
		colTypes[i] = ct.ScanType().String()
		if ct.ScanType().Kind() == reflect.Interface {
			if v, ok := s.Typemap.GetGoType(ct.DatabaseTypeName()); ok {
				colTypes[i] = v
			} else {
				return nil, nil, fmt.Errorf("%s: %s", errColTypes, ct.DatabaseTypeName())
			}
		}
	}

	return cols, colTypes, nil
}

// GenerateStruct .
func (s *Stst) GenerateStruct(name string, cols [][3]string) (*jen.Statement, error) {
	// Use []jen.Code instead of []*jen.Statement to pass it to jen.Structs()
	ms := make([]jen.Code, len(cols))
	for i, c := range cols {
		ms[i] = jen.Id(strcase.ToCamel(c[0])).Add(jen.Qual(c[1], c[2]))
	}

	st := jen.Type().Id(name).Struct(ms...)
	return st, nil
}

// GenerateGetScanDestsFunc .
func (s *Stst) GenerateGetScanDestsFunc(structName string, fieldNames []string) (*jen.Statement, error) {
	recn := "x"
	rec := jen.Id(recn).Op("*").Id(structName)

	fn := "GetScanDests"
	sig := jen.Func().Params(rec).Id(fn).Params().Index().Interface()

	fields := make([]jen.Code, len(fieldNames))
	for i, n := range fieldNames {
		fields[i] = jen.Op("&").Id(recn).Dot(n)
	}

	ret := jen.Return(
		jen.Index().Interface().Values(fields...),
	)

	res := sig.Block(ret)

	return res, nil
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
