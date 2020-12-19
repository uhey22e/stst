CREATE TABLE public.basic_types (
    bigint_col bigint NOT NULL PRIMARY KEY
    , text_col text NOT NULL
);

INSERT INTO public.basic_types (bigint_col , text_col)
    VALUES (1 , 'Example text') , (2 , 'A long time ago in a galaxy far, far away...');

CREATE TABLE public.nullable (
    bigint_col bigserial PRIMARY KEY
    , nullable_col text
);

