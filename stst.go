package stst

import (
	"database/sql"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

// Stst .
type Stst struct {
	DB      *sql.DB
	Typemap Typemap
}

// ColInfo is a concrete representation of a column and a member of struct.
type ColInfo struct {
	Name        string
	GoTypeName  string
	PackagePath string
}

// DBConf contains configurations connecting to database.
type DBConf struct {
	Host     string `env:"DB_HOST" envDefault:"127.0.0.1"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	DBName   string `env:"DB_DBNAME" envDefault:"postgres"`
	Username string `env:"DB_USERNAME" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"postgres"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}

// Typemap is an interface
type Typemap interface {
	// GetGoType returns a golang type name mapped from the database type name
	GetGoType(databaseTypeName string) (string, bool)
}

type MemberCustomizer func(*jen.Statement, *ColInfo)

var (
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

// GenerateStruct assembles the struct from ColInfo.
func (s *Stst) GenerateStruct(name string, cols []ColInfo, memberCustomizer ...MemberCustomizer) (*jen.Statement, error) {
	st := jen.Type().Id(name).StructFunc(func(g *jen.Group) {
		for _, c := range cols {
			n := strcase.ToCamel(c.Name)
			m := jen.Id(n).Qual(c.PackagePath, c.GoTypeName)
			for _, mc := range memberCustomizer {
				mc(m, &c)
			}
			g.Add(m)
		}
	})
	return st, nil
}

// GenerateGetScanDestsFunc assembles the GetScanDests() function from ColInfo.
func (s *Stst) GenerateGetScanDestsFunc(structName string, cols []ColInfo) (*jen.Statement, error) {
	recn := "x"
	rec := jen.Id(recn).Op("*").Id(structName)

	fn := "GetScanDests"
	rettype := jen.Index().Interface() // []interface{}
	sig := jen.Func().Params(rec).Id(fn).Params().Add(rettype)

	ret := jen.Return(
		jen.Index().Interface().ValuesFunc(func(g *jen.Group) {
			for _, c := range cols {
				g.Op("&").Id(recn).Dot(strcase.ToCamel(c.Name))
			}
		}),
	)

	res := sig.Block(ret)

	return res, nil
}

// GenerateQueryVar assembles the variable definition of the query.
func (s *Stst) GenerateQueryVar(name, query string) (*jen.Statement, error) {
	q := strings.Trim(query, "\n")
	vn := name + "Query"
	return jen.Var().Id(vn).Op("=").Id(fmt.Sprintf("`\n%s`", q)), nil
}

// Package .
func (s *Stst) Package(w io.Writer, name string, codes []jen.Code, pkgComments []string) error {
	f := jen.NewFile(name)
	if pkgComments != nil {
		for _, c := range pkgComments {
			f.PackageComment(c)
		}
	}

	for _, c := range codes {
		f.Add(c)
	}
	if err := f.Render(w); err != nil {
		return err
	}
	return nil
}

// AppendColNameTag is an option for GenerateStruct()
func AppendColNameTag(key string) MemberCustomizer {
	return func(s *jen.Statement, c *ColInfo) {
		s.Tag(map[string]string{
			key: c.Name,
		})
	}
}

func trimSemicolon(q string) string {
	r := regexp.MustCompile(`;[\s]*$`)
	q2 := r.ReplaceAllString(q, "")
	return q2
}
