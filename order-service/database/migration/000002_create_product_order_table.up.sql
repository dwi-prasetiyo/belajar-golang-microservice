CREATE TABLE IF NOT EXISTS product_orders (
    id SERIAL PRIMARY KEY,
    order_id VARCHAR(21) NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    price INT NOT NULL
)