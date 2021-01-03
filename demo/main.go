package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/uhey22e/stst/demo/models"
)

var (
	db *sql.DB
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getPopularActor() {
	r, err := db.Query(models.PopularActorQuery)
	handleError(err)
	defer r.Close()

	for r.Next() {
		var x models.PopularActor
		r.Scan(x.GetScanDests()...)
		log.Printf("%+v", x)
	}
}

func getHardworkingStaff() {
	r, err := db.Query(models.HardworkingStaffQuery)
	handleError(err)
	defer r.Close()

	for r.Next() {
		var x models.HardworkingStaff
		r.Scan(x.GetScanDests()...)
		log.Printf("%+v", x)
	}
}

func main() {
	var err error
	dsn := "postgresql://postgres:postgres@localhost:15432/dvdrental?sslmode=disable"
	db, err = sql.Open("postgres", dsn)
	handleError(err)

	err = db.Ping()
	handleError(err)

	getPopularActor()
	getHardworkingStaff()
}
