CREATE TABLE IF NOT EXISTS products (
	id BYTEA PRIMARY KEY,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP,
	deleted_at TIMESTAMP,

	sku varchar(25) NOT NULL,
	barcode varchar(50) NOT NULL,
	name varchar(25) NOT NULL,
	description varchar(100) NOT NULL,
	image BYTEA,
	amount NUMERIC(12,2) NOT NULL
);
