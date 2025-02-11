-- +goose Up
CREATE TABLE merch (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    price INT NOT NULL CHECK (price > 0)
);

-- Дефолтные значения
INSERT INTO merch (name, price) VALUES
    ('t-shirt', 80), ('cup', 20), ('book', 50), ('pen', 10),
    ('powerbank', 200), ('hoody', 300), ('umbrella', 200),
    ('socks', 10), ('wallet', 50), ('pink-hoody', 500);

-- +goose Down
DROP TABLE merch;
