-- +goose Up
CREATE TABLE purchases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    merch_id INT NOT NULL REFERENCES merch(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE purchases;

