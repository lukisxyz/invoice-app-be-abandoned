CREATE TABLE IF NOT EXISTS category_products (
    category_id BYTEA,
    product_id BYTEA,
    PRIMARY KEY (category_id, product_id),
    FOREIGN KEY (category_id) REFERENCES categories(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);