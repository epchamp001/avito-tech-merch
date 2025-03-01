-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    balance INT NOT NULL DEFAULT 1000,
    created_at TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE users;
