-- +goose Up
ALTER TABLE users DROP CONSTRAINT users_email_key;

CREATE UNIQUE INDEX users_email_lower_key
    ON users (LOWER(email));

-- +goose Down
DROP INDEX users_email_lower_key;

ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);