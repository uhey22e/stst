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
	)
	flag.Parse()

	dsn := "postgresql://postgres:postgres@localhost:15432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	handleError(err)

	s := stst.New(db)

	q, err := ioutil.ReadFile(*sqlFile)
	handleError(err)

	_, colTypes, err := s.GetMeta(string(q))
	handleError(err)

	members := make([][2]string, len(colTypes))
	for i, c := range colTypes {
		members[i] = [2]string{
			c.Name(),
			c.ScanType().Name(),
		}
	}

	st, err := s.GenerateStruct(members)
	handleError(err)

	b := &bytes.Buffer{}
	s.Package(b, "simple", []jen.Code{st})

	fmt.Printf("%s", b.String())
}
