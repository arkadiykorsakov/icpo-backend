-- +goose Up
CREATE TABLE hair_colors (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX hair_colors_name_ci
    ON hair_colors (LOWER(name));

-- +goose Down
DROP TABLE hair_colors;
