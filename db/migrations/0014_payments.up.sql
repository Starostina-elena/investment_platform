CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY,
    external_id TEXT NOT NULL UNIQUE,
    amount NUMERIC(34, 2) NOT NULL,
    user_id INT NOT NULL REFERENCES users(id),
    entity_type TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_payments_external_id ON payments(external_id);
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);
