-- +goose Up
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID REFERENCES users(id) ON DELETE SET NULL,
    receiver_id UUID REFERENCES users(id) ON DELETE CASCADE,
    amount INT NOT NULL CHECK (amount > 0),
    type TEXT NOT NULL CHECK (type IN ('initial', 'transfer', 'purchase')),
    created_at TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE transactions;
