package stst

import (
	"bytes"
	"database/sql"
	"fmt"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/google/go-cmp/cmp"
	_ "github.com/lib/pq"
)

func testConnectDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := "postgresql://postgres:postgres@localhost:15432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func testLoadQuery(t *testing.T, filename string) string {
	t.Helper()

	q, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return string(q)
}

func formatGoCode(t *testing.T, code string) string {
	t.Helper()

	buf := bytes.NewBufferString(code)
	fmt, err := format.Source(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	res := string(fmt)

	return res
}

func TestCountWrap(t *testing.T) {
	q := `select * from basic_types limit 10;`
	wrapped := fmt.Sprintf(`select count(*) from (%s) x1;`, trimSemicolon(q))
	t.Log(wrapped)
}

func TestStst_GetMeta(t *testing.T) {
	type fields struct {
		db *sql.DB
	}

	s := NewPsql(testConnectDB(t))

	tests := []struct {
		name    string
		sqlFile string
		want    []ColInfo
		wantErr bool
	}{
		{
			"BasicTypes",
			filepath.Join("testdata", "basic_types.sql"),
			[]ColInfo{
				{
					Name:        "bigint_col",
					GoTypeName:  "int64",
					PackagePath: "",
				},
				{
					Name:        "text_col",
					GoTypeName:  "string",
					PackagePath: "",
				},
			},
			false,
		},
		{
			"ComplexTypes",
			filepath.Join("testdata", "complex_types.sql"),
			[]ColInfo{
				{
					Name:        "bigint_col",
					GoTypeName:  "int64",
					PackagePath: "",
				},
				{
					Name:        "double_precision_col",
					GoTypeName:  "float64",
					PackagePath: "",
				},
				{
					Name:        "timestamp_col",
					GoTypeName:  "Time",
					PackagePath: "time",
				},
				{
					Name:        "numeric_col",
					GoTypeName:  "float64",
					PackagePath: "",
				},
			},
			false,
		},
		{
			"Nullable",
			filepath.Join("testdata", "nullable.sql"),
			[]ColInfo{
				{
					Name:        "bigint_col",
					GoTypeName:  "int64",
					PackagePath: "",
				},
				{
					Name:        "nullable_col",
					GoTypeName:  "string",
					PackagePath: "",
				},
			},
			false,
		},
		{
			"InvalidQuery",
			filepath.Join("testdata", "invalid.sql"),
			[]ColInfo{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := testLoadQuery(t, tt.sqlFile)
			cols, err := s.GetMeta(q)

			if (err != nil) != tt.wantErr {
				t.Errorf("Stst.GetMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil {
				t.Logf("Expected error: %v", err)
				return
			}

			if !cmp.Equal(cols, tt.want) {
				t.Errorf("Stst.GetMeta() got = %v, want %v", cols, tt.want)
			}
		})
	}
}

func TestStst_GenerateStruct(t *testing.T) {
	type fields struct {
		db *sql.DB
	}

	s := NewPsql(nil)
	tests := []struct {
		name    string
		sname   string
		cols    []ColInfo
		want    string
		wantErr bool
	}{
		{
			"Case1",
			"Demo",
			[]ColInfo{
				{
					Name:        "col1",
					GoTypeName:  "int64",
					PackagePath: "",
				},
				{
					Name:        "col2",
					GoTypeName:  "string",
					PackagePath: "",
				},
				{
					Name:        "col3",
					GoTypeName:  "float64",
					PackagePath: "",
				},
			},
			formatGoCode(t, `
				package models

				type Demo struct {
					Col1 int64
					Col2 string
					Col3 float64
				}
			`),
			false,
		},
		{
			"Case2",
			"Demo",
			[]ColInfo{
				{
					Name:        "bigint_col",
					GoTypeName:  "int64",
					PackagePath: "",
				},
				{
					Name:        "timestamp_col",
					GoTypeName:  "Time",
					PackagePath: "time",
				},
			},
			formatGoCode(t, `
				package models

				import "time"

				type Demo struct {
					BigintCol    int64
					TimestampCol time.Time
				}
			`),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GenerateStruct(tt.sname, tt.cols)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stst.GenerateStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			f := jen.NewFile("models")
			f.Add(got)
			gotfmt := f.GoString()

			if gotfmt != tt.want {
				t.Errorf("Stst.GenerateStruct() got =\n%v\nwant =\n%v", gotfmt, tt.want)
			}
		})
	}
}

func TestStst_GenerateGetScanDestsFunc(t *testing.T) {
	s := NewPsql(nil)
	cols := []ColInfo{
		{
			Name: "Col1",
		},
		{
			Name: "Col2",
		},
		{
			Name: "Col3",
		},
	}

	want := formatGoCode(t, `
		package models

		func (x *Demo) GetScanDests() []interface{} {
			return []interface{}{&x.Col1, &x.Col2, &x.Col3}
		}`)

	res, err := s.GenerateGetScanDestsFunc("Demo", cols)
	if err != nil {
		t.Fatal(err)
	}

	f := jen.NewFile("models")
	f.Add(res)
	got := f.GoString()

	if got != want {
		t.Errorf("Stst.GenerateGetScanDestsFunc() got =\n%v\nwant =\n%v", got, want)
	}
}

func TestStst_Package(t *testing.T) {
	s := NewPsql(nil)
	w := &bytes.Buffer{}

	codes := []jen.Code{
		jen.Type().Id("TestStruct").Struct(
			jen.Id("Member1").Int(),
			jen.Id("Member2").Float64(),
		),
		// def (ts TestStruct) Method() string {}
		jen.Func().Params(jen.Id("ts").Id("*TestStruct")).Id("Method").Params().String().Block(
			jen.Return(jen.Id("ts.Member1")),
		),
	}

	if err := s.Package(w, "test", codes, nil); err != nil {
		t.Fatalf("Stst.Package() error = %v", err)
	} else {
		t.Logf("\n%s", w.String())
	}
}

func Test_trimSemicolon(t *testing.T) {
	tests := []struct {
		test string
		want string
	}{
		{
			"select * from table;",
			"select * from table",
		},
		{
			`select
				col01,
				col02
			from table01
			;
			`,
			`select
				col01,
				col02
			from table01
			`,
		},
	}
	for i, tt := range tests {
		t.Run("Case"+strconv.Itoa(i), func(t *testing.T) {
			if got := trimSemicolon(tt.test); got != tt.want {
				t.Errorf("trimSemicolon() = %v, want %v", got, tt.want)
			}
		})
	}
}
