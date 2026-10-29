-- +goose Up
CREATE TABLE cities (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    country_id BIGINT NOT NULL REFERENCES countries (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- один город с данным именем (без учёта регистра) на страну
CREATE UNIQUE INDEX cities_name_ci
    ON cities (LOWER(name), country_id);

-- фильтр ?country_id= и проверка FK при удалении страны
CREATE INDEX cities_country_id_idx
    ON cities (country_id);

-- +goose Down
DROP TABLE cities;