package models

type Demo struct {
	BigintCol int64
	TextCol   string
}

const DemoQuery = `
SELECT
    bigint_col,
    text_col
FROM
    basic_types
LIMIT 10;
`

var DemoCountQuery = `
SELECT
    count(*)
FROM (
    SELECT
        bigint_col,
        text_col
    FROM
        basic_types
    LIMIT 10) x_1;
`

func (x *Demo) GetScanDests() []interface{} {
	return []interface{}{
		&x.BigintCol,
		&x.TextCol,
	}
}
