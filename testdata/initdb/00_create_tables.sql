CREATE TABLE public.simple (
    bigint_col bigserial NOT NULL PRIMARY KEY
    , text_col text NOT NULL
    , timestamp_col timestamp NOT NULL
);

CREATE TABLE public.nullable (
    bigint_col bigserial PRIMARY KEY
    , nullable_col text
);

