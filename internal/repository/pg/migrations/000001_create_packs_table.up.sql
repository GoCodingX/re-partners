CREATE TABLE packs
(
    id         UUID PRIMARY KEY NOT NULL,
    size       INTEGER UNIQUE   NOT NULL,
    created_at TIMESTAMPTZ      NOT NULL NOT NULL,
    updated_at TIMESTAMPTZ      NOT NULL NOT NULL
);
