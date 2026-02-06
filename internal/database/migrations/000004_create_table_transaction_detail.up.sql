CREATE TABLE transaction_detail (
    id SERIAL PRIMARY KEY,
    transaction_id INT REFERENCES transaction(id) ON DELETE CASCADE,
    product_id INT REFERENCES product(id),
    quantity INT NOT NULL,
    subtotal INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);