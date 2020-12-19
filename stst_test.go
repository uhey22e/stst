package stst

import (
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
	tests := []struct {
		name    string
		fields  fields
		want    *jen.Statement
		wantErr bool
	}{
		{
			"Case",
			fields{nil},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stst{
				db: tt.fields.db,
			}
			got, err := s.GenerateStruct()
			if (err != nil) != tt.wantErr {
				t.Errorf("Stst.GenerateStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("\n%#v\n", got)
		})
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
