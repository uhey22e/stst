package stst

import (
	"database/sql"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

// Stst .
type Stst struct {
	DB      *sql.DB
	Typemap Typemap
}

type ColInfo struct {
	Name        string
	GoTypeName  string
	PackagePath string
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
func (s *Stst) GetMeta(query string) ([]ColInfo, error) {
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errCols, err)
	}

	cts, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errColTypes, err)
	}

	ms := make([]ColInfo, len(cols))
	for i, ct := range cts {
		ms[i].Name = ct.Name()
		ms[i].GoTypeName = ct.ScanType().String()
		ms[i].PackagePath = ""

		if ct.ScanType().Kind() == reflect.Interface {
			if v, ok := s.Typemap.GetGoType(ct.DatabaseTypeName()); ok {
				ms[i].GoTypeName = v
			} else {
				return nil, fmt.Errorf("%s: %s", errColTypes, ct.DatabaseTypeName())
			}
		}

		s := strings.Split(ms[i].GoTypeName, ".")
		if len(s) == 2 {
			ms[i].PackagePath = s[0]
			ms[i].GoTypeName = s[1]
		}
	}

	return ms, nil
}

// GenerateStruct .
func (s *Stst) GenerateStruct(name string, cols []ColInfo) (*jen.Statement, error) {
	// Use []jen.Code instead of []*jen.Statement to pass it to jen.Structs()
	ms := make([]jen.Code, len(cols))
	for i, c := range cols {
		n := strcase.ToCamel(c.Name)
		ms[i] = jen.Id(n).Add(jen.Qual(c.PackagePath, c.GoTypeName))
	}

	st := jen.Type().Id(name).Struct(ms...)
	return st, nil
}

// GenerateGetScanDestsFunc .
func (s *Stst) GenerateGetScanDestsFunc(structName string, cols []ColInfo) (*jen.Statement, error) {
	recn := "x"
	rec := jen.Id(recn).Op("*").Id(structName)

	fn := "GetScanDests"
	rettype := jen.Index().Interface() // []interface{}
	sig := jen.Func().Params(rec).Id(fn).Params().Add(rettype)

	fields := make([]jen.Code, len(cols))
	for i, c := range cols {
		fields[i] = jen.Op("&").Id(recn).Dot(strcase.ToCamel(c.Name))
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
