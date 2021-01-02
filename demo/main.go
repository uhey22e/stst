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

	r1, err := db.Query(models.DemoCountQuery)
	handleError(err)
	defer r1.Close()

	var count int64
	if v := r1.Next(); !v {
		log.Printf("Failed to get count")
	}
	r1.Scan(&count)
	log.Printf("Count %d", count)

	if count == 0 {
		return
	}

	r2, err := db.Query(models.DemoQuery)
	handleError(err)
	defer r2.Close()

	result := make([]models.Demo, 0, count)
	for r2.Next() {
		var x models.Demo
		r2.Scan(x.GetScanDests()...)
		result = append(result, x)
	}

	log.Printf("%+v\n", result)
}
