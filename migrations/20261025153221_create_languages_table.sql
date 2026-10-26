-- +goose Up
CREATE TABLE languages (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX languages_name_ci
    ON languages (lower(name));

-- +goose Down
DROP TABLE languages;
