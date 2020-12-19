package stst

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/google/go-cmp/cmp"
	_ "github.com/lib/pq"
)

func TestStst_GetMeta(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		query string
	}

	s := &Stst{
		db: testConnectDB(t),
	}

	tests := []struct {
		name         string
		sqlFile      string
		wantCols     []string
		wantColTypes []*sql.ColumnType
		wantErr      bool
	}{
		{
			"Simple",
			filepath.Join("testdata", "simple.sql"),
			[]string{"bigint_col", "text_col", "timestamp_col"},
			[]*sql.ColumnType{},
			false,
		},
		{
			"InvalidQuery",
			filepath.Join("testdata", "invalid.sql"),
			[]string{},
			[]*sql.ColumnType{},
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
			}
			if err != nil {
				t.Logf("Expected error: %v", err)
				return
			}

			if !cmp.Equal(cols, tt.wantCols) {
				t.Errorf("Stst.GetMeta() got = %v, want %v", cols, tt.wantCols)
			}
			// TODO: Test database type
			for _, ct := range colTypes {
				t.Log(ct)
			}
		})
	}
}

func TestStst_GenerateStruct(t *testing.T) {
	type fields struct {
		db *sql.DB
	}

	s := &Stst{nil}
	tests := []struct {
		name    string
		cols    [][2]string
		want    *jen.Statement
		wantErr bool
	}{
		{
			"Case",
			[][2]string{
				{"bigint_col", "int64"},
				{"timestamp_col", "time.Time"},
			},
			jen.Type().Id("Foo").Struct(
				jen.Id("BigintCol").Int64(),
				jen.Id("TimestampCol").Add(jen.Id("time.Time")),
			),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GenerateStruct(tt.cols)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stst.GenerateStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !cmp.Equal(got.GoString(), tt.want.GoString()) {
				t.Errorf("Stst.GenerateStruct() got = %v, want %v", got, tt.want)
			} else {
				t.Logf("\n%#v\n", got)
			}
		})
	}
}

func TestStst_Package(t *testing.T) {
	s := &Stst{nil}
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
