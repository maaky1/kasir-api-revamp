CREATE TABLE transaction (
    id   SERIAL PRIMARY KEY,
    total_amount INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);