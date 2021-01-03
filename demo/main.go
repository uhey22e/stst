package main

import (
	"database/sql"
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

	r, err := db.Query(models.DemoQuery)
	handleError(err)
	defer r.Close()

	for r.Next() {
		var x models.Demo
		r.Scan(x.GetScanDests()...)
		log.Printf("%+v", x)
	}
}
