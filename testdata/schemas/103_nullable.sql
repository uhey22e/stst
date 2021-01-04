DROP TABLE IF EXISTS public.nullable;

CREATE TABLE public.nullable (
    bigint_col bigserial PRIMARY KEY,
    nullable_col text
);

INSERT INTO public.nullable (bigint_col, nullable_col)
    VALUES (1, NULL), (2, 'NOT NULL');

