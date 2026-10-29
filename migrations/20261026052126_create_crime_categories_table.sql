-- +goose Up
CREATE TABLE crime_categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX crime_categories_name_ci
    ON crime_categories (LOWER(name));

-- +goose Down
DROP TABLE crime_categories;
