CREATE TABLE IF NOT EXISTS products (
	id BYTEA PRIMARY KEY,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP,
	deleted_at TIMESTAMP,

	sku varchar(25) NOT NULL UNIQUE,
	name varchar(100) NOT NULL,
	description Text NOT NULL,
	image BYTEA,
	amount NUMERIC(12,2) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sku ON products(sku);