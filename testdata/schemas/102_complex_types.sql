DROP TABLE IF EXISTS public.complex_types;

CREATE TABLE public.complex_types (
    bigint_col bigint NOT NULL PRIMARY KEY,
    double_precision_col double precision NOT NULL,
    timestamp_col timestamptz NOT NULL,
    numeric_col numeric(8, 3) NOT NULL
);

INSERT INTO public.complex_types (bigint_col, double_precision_col, timestamp_col, numeric_col)
    VALUES (1, 123.45, '2020-10-01 15:00:00', 123.45), (2, 2345.678, '2020-10-01 15:00:00', 2345.678);

