-- +goose Up
CREATE TABLE countries (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX countries_name_ci
    ON countries (lower(name));

-- +goose Down
DROP TABLE countries;
