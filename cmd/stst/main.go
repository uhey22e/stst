package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/dave/jennifer/jen"
	_ "github.com/lib/pq"
	"github.com/uhey22e/stst"
)

var comments = []string{
	"Code generated by stst; DO NOT EDIT.",
}

var (
	version = "dev"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var (
		showVersion   = flag.Bool("v", false, "Show version.")
		sqlFile       = flag.String("i", "", "Input SQL file.")
		name          = flag.String("p", "models", "Output package name.")
		sname         = flag.String("n", "", "Output struct type name.")
		outFile       = flag.String("o", "", "Output file name. Output to stdout if empty.")
		appendBoilTag = flag.Bool("b", false, "Append tag to the struct for SQLBoiler.")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("stst %s\n", version)
		os.Exit(0)
	}

	dbconf := stst.DBConf{}
	err := env.Parse(&dbconf)
	handleError(err)
	conninfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", dbconf.Host, dbconf.Port, dbconf.Username, dbconf.Password, dbconf.DBName, dbconf.SSLMode)
	db, err := sql.Open("postgres", conninfo)
	handleError(err)

	// Supports PostgreSQL only
	s := stst.NewPsql(db)

	q, err := ioutil.ReadFile(*sqlFile)
	handleError(err)

	cols, err := s.GetMeta(string(q))
	handleError(err)

	cs := []stst.MemberCustomizer{}
	if *appendBoilTag {
		cs = append(cs, stst.AppendColNameTag("boil"))
	}
	st, err := s.GenerateStruct(*sname, cols, cs...)
	handleError(err)

	qv, err := s.GenerateQueryVar(*sname, string(q))
	handleError(err)

	f, err := s.GenerateGetScanDestsFunc(*sname, cols)
	handleError(err)

	var dest io.Writer
	if *outFile != "" {
		f, err := os.Create(*outFile)
		handleError(err)
		defer f.Close()
		dest = f
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
