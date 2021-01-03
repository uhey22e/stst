package main

import (
	"database/sql"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/dave/jennifer/jen"
	_ "github.com/lib/pq"
	"github.com/uhey22e/stst"
)

var comments = []string{
	"DO NOT EDIT THIS FILE MANUALLY.",
	"This file is generated by stst.",
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var (
		sqlFile = flag.String("i", "", "Input SQL file")
		name    = flag.String("p", "models", "Output package name.")
		outFile = flag.String("o", "", "Output file name. Output to stdout if empty.")
	)
	flag.Parse()

	dsn := "postgresql://postgres:postgres@localhost:15432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	handleError(err)

	// Supports PostgreSQL only
	s := stst.NewPsql(db)

	q, err := ioutil.ReadFile(*sqlFile)
	handleError(err)

	cols, err := s.GetMeta(string(q))
	handleError(err)

	st, err := s.GenerateStruct("Demo", cols)
	handleError(err)

	qv, err := s.GenerateQueryVar("Demo", string(q))
	handleError(err)

	f, err := s.GenerateGetScanDestsFunc("Demo", cols)
	handleError(err)

	var dest io.Writer
	if *outFile != "" {
		dest, err := os.Create(*outFile)
		handleError(err)
		defer dest.Close()
	} else {
		dest = os.Stdout
	}

	err = s.Package(dest, *name, []jen.Code{
		qv,
		st,
		f,
	}, comments)
	handleError(err)

}
