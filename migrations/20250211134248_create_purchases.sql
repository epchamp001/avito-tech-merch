-- +goose Up
CREATE TABLE purchases (
    id SERIAL PRIMARY KEY ,
    user_id int NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    merch_id INT NOT NULL REFERENCES merch(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE purchases;

