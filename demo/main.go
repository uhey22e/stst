package main

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
	"github.com/uhey22e/stst/demo/models"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	dsn := "postgresql://postgres:postgres@localhost:15432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	handleError(err)

	q, err := ioutil.ReadFile("testdata/basic_types.sql")
	handleError(err)

	rows, err := db.Query(string(q))
	handleError(err)

	for rows.Next() {
		var m models.Foo
		rows.Scan(&m.BigintCol, &m.TextCol)
		log.Printf("%+v\n", m)
	}
}
