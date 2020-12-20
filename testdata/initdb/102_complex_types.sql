CREATE TABLE public.complex_types (
    bigint_col bigint NOT NULL PRIMARY KEY,
    numeric_precision_col double precision NOT NULL,
    timestamp_col timestamptz NOT NULL,
    decimal_col numeric(8, 3) NOT NULL
);

INSERT INTO public.complex_types (bigint_col, numeric_precision_col, timestamp_col, decimal_col)
    VALUES (1, 123.45, '2020-10-01 15:00:00', 123.45), (2, 2345.678, '2020-10-01 15:00:00', 2345.678);

