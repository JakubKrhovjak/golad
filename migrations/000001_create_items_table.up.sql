CREATE TABLE IF NOT EXISTS items (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    name VARCHAR(255) NOT NULL
    description TEXT
);

CREATE INDEX IF NOT EXISTS idx_items_name ON items(name);
