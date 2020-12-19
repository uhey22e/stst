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

	dsn := "postgresql://postgres:postgres@localhost:15432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	sqlFile := filepath.Join("testdata", "simple.sql")
	q, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		t.Fatal(err)
	}
	query := string(q)

	tests := []struct {
		name         string
		fields       fields
		args         args
		wantCols     []string
		wantColTypes []*sql.ColumnType
		wantErr      bool
	}{
		{
			"Case",
			fields{db},
			args{query},
			[]string{"bigint_col", "text_col", "timestamp_col"},
			[]*sql.ColumnType{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stst{
				db: tt.fields.db,
			}
			cols, colTypes, err := s.GetMeta(tt.args.query)

			if (err != nil) != tt.wantErr {
				t.Errorf("Stst.GetMeta() error = %v, wantErr %v", err, tt.wantErr)
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
