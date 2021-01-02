package stst

import (
	"bytes"
	"database/sql"
	"fmt"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/google/go-cmp/cmp"
	_ "github.com/lib/pq"
)

func TestCountWrap(t *testing.T) {
	q := `select * from basic_types limit 10`
	wrapped := fmt.Sprintf(`select count(*) from (%s) x1;`, q)
	t.Log(wrapped)
}

func TestStst_GetMeta(t *testing.T) {
	type fields struct {
		db *sql.DB
	}

	s := NewPsql(testConnectDB(t))

	tests := []struct {
		name         string
		sqlFile      string
		wantCols     []string
		wantColTypes []string
		wantErr      bool
	}{
		{
			"BasicTypes",
			filepath.Join("testdata", "basic_types.sql"),
			[]string{"bigint_col", "text_col"},
			[]string{"int64", "string"},
			false,
		},
		{
			"ComplexTypes",
			filepath.Join("testdata", "complex_types.sql"),
			[]string{"bigint_col", "double_precision_col", "timestamp_col", "numeric_col"},
			[]string{"int64", "float64", "time.Time", "float64"},
			false,
		},
		{
			"Nullable",
			filepath.Join("testdata", "nullable.sql"),
			[]string{"bigint_col", "nullable_col"},
			[]string{"int64", "string"},
			false,
		},
		{
			"InvalidQuery",
			filepath.Join("testdata", "invalid.sql"),
			[]string{},
			[]string{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := testLoadQuery(t, tt.sqlFile)
			cols, colTypes, err := s.GetMeta(q)

			if (err != nil) != tt.wantErr {
				t.Errorf("Stst.GetMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil {
				t.Logf("Expected error: %v", err)
				return
			}

			if !cmp.Equal(cols, tt.wantCols) {
				t.Errorf("Stst.GetMeta() cols got = %v, want %v", cols, tt.wantCols)
			}

			if !cmp.Equal(colTypes, tt.wantColTypes) {
				t.Errorf("Stst.GetMeta() colTypes got = %v, want %v", colTypes, tt.wantColTypes)
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
		cols    [][3]string
		want    string
		wantErr bool
	}{
		{
			"Case1",
			"Demo",
			[][3]string{
				{"col1", "", "int64"},
				{"col2", "", "string"},
				{"col3", "", "float64"},
				{"col4", "", "string"},
			},
			formatGoCode(t, `
				package models

				type Demo struct {
					Col1 int64
					Col2 string
					Col3 float64
					Col4 string
				}
			`),
			false,
		},
		{
			"Case2",
			"Demo",
			[][3]string{
				{"bigint_col", "", "int64"},
				{"timestamp_col", "time", "Time"},
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
	cols := []string{"Col1", "Col2", "Col3"}

	want := formatGoCode(t, `
		package models

		func (x *Demo) GetScanDests() []interface{} {
			return []interface{}{&x.Col1, &x.Col2, &x.Col3}
		}
	`)

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

	if err := s.Package(w, "test", codes); err != nil {
		t.Fatalf("Stst.Package() error = %v", err)
	} else {
		t.Logf("\n%s", w.String())
	}
}

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
