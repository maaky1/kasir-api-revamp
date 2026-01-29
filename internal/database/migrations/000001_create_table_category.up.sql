CREATE TABLE category (
    id   SERIAL PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	description TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO category (name, description, created_at)
VALUES ('Uncategorized', 'Default category', NOW());