# stst

**Notice** : This project is still a work in progress.

Stst is an code generator that creates structs from SQL queries.


## How to use

```
$ cat sql/example.sql

SELECT
    foo_id
    , name
    , bar_date
FROM
    awesome_table
LIMIT 1;


$ stst -i testdata/example.sql -p models

package models

type Example struct {
    FooID   int64
    Name    string
    BarDate Time
}
...
```

```go
// your_app.go

rows, err := db.Query(q)
for rows.Next() {
    var m models.Example
    err := m.Scan(rows)
    fmt.Printf("%#v", m)
}
```

To be written...

## Options

To be written...

## Usage scenarios

To be written...
