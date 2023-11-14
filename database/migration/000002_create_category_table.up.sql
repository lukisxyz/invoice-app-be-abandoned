CREATE TABLE IF NOT EXISTS categories (
	id BYTEA PRIMARY KEY,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP,
	deleted_at TIMESTAMP,

	name varchar(50) NOT NULL,
	description varchar(100) NOT NULL
);
