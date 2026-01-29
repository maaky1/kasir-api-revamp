CREATE TABLE product (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL DEFAULT 1,
    name TEXT NOT NULL,
    price INT NOT NULL,
    stock INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_product_category
        FOREIGN KEY (category_id)
        REFERENCES category(id)
        ON DELETE SET DEFAULT
);