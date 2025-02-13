-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    balance INT NOT NULL DEFAULT 1000,
    created_at TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE users;
