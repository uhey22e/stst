package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
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

	members := make([][3]string, len(colTypes))
	for i := 0; i < len(cols); i++ {
		s := strings.Split(colTypes[i], ".")
		if len(s) == 2 {
			members[i] = [3]string{
				cols[i],
				s[0],
				s[1],
			}
		} else {
			members[i] = [3]string{
				cols[i],
				"",
				colTypes[i],
			}
		}
	}

	st, err := s.GenerateStruct("Demo", members)
	handleError(err)

	qv := jen.Const().Id("DemoQuery").Op("=").Id(fmt.Sprintf("`\n%s`", q))

	m2 := make([]string, len(members))
	for i, m := range members {
		m2[i] = strcase.ToCamel(m[0])
	}

	f, err := s.GenerateGetScanDestsFunc("Demo", m2)
	handleError(err)

	b := &bytes.Buffer{}
	s.Package(b, *name, []jen.Code{
		qv,
		st,
		f,
	})

	if *outFile != "" {
		err := ioutil.WriteFile(*outFile, b.Bytes(), 0644)
		handleError(err)
	} else {
		fmt.Printf("%s", b.String())
	}
}
