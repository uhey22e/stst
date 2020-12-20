package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/dave/jennifer/jen"
	_ "github.com/lib/pq"
	"github.com/uhey22e/stst"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var (
		sqlFile = flag.String("i", "", "Input SQL file")
		name    = flag.String("p", "models", "Output package name")
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

	cols, colTypes, err := s.GetMeta(string(q))
	handleError(err)

	members := make([][2]string, len(colTypes))
	for i := 0; i < len(cols); i++ {
		members[i] = [2]string{
			cols[i],
			colTypes[i],
		}
	}

	st, err := s.GenerateStruct(members)
	handleError(err)

	b := &bytes.Buffer{}
	s.Package(b, *name, []jen.Code{st})

	if *outFile != "" {
		err := ioutil.WriteFile(*outFile, b.Bytes(), 0644)
		handleError(err)
	} else {
		fmt.Printf("%s", b.String())
	}
}
