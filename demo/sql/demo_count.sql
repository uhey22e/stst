SELECT
    count(*)
FROM (
    SELECT
        bigint_col,
        text_col
    FROM
        basic_types
    LIMIT 10) x_1;

