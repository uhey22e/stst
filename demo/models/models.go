//go:generate env DB_DBNAME=dvdrental DB_PORT=15432 go run github.com/uhey22e/stst/cmd/stst -i ../sql/popular_actor.sql -o ./popular_actor.go -n PopularActor
//go:generate env DB_DBNAME=dvdrental DB_PORT=15432 go run github.com/uhey22e/stst/cmd/stst -i ../sql/hardworking_staff.sql -o ./hardworking_staff.go -n HardworkingStaff

package models
