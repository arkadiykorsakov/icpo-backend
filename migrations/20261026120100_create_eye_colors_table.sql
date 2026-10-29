-- +goose Up
CREATE TABLE eye_colors (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX eye_colors_name_ci
    ON eye_colors (LOWER(name));

-- +goose Down
DROP TABLE eye_colors;
