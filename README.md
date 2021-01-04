# stst

**Notice** : This project is still a work in progress.

Stst is an code generator which generates Go structs from SELECT statements.


## How to use

Provide a SELECT statement.

```sql
-- sql/example.sql
SELECT
    foo_id
    , bar_name
FROM
    awesome_table
LIMIT 1;
```

Generate a Go source code from the sql file.

```
$ go run github.com/uhey22e/stst/cmd/stst -i sql/example.sql -p models -n Example -o models/example.go
$ less models/example.go

package models

type Example struct {
    FooID   int64
    BarName string
}
...
```

Use the generated struct and helper functions for the `database/sql` package in your application.

```go
// your_app.go

rows, err := db.Query(models.ExampleQuery)
for rows.Next() {
    var m models.Example
    err := rows.Scan(m.GetScanDests()...)
    fmt.Printf("%#v", m)
}
```

For more examples, please see [demo project](./demo).


## Options

To be written...


## Tests

A runnning PostgreSQL server is required for tests.

Run local PostgreSQL server using docker:

```
$ cd testdb
$ docker-compose up -d
```

Run tests:

```
$ go test -v
```

