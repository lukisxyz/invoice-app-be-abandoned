CREATE TABLE IF NOT EXISTS inventories (
    id BYTEA PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    quantity INTEGER NOT NULL
);

ALTER TABLE products
ADD COLUMN inventory_id BYTEA UNIQUE;

ALTER TABLE products
ADD CONSTRAINT fk_inventory
FOREIGN KEY (inventory_id)
REFERENCES inventories(id);