-- +goose Up
CREATE TABLE eye_colors (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX eye_colors_name_ci
    ON eye_colors (LOWER(name));

-- +goose Down
DROP TABLE eye_colors;
