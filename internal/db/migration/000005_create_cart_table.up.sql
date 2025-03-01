CREATE TABLE cart (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id INT NOT NULL,
    seller_id INT NOT NULL,
    product_id INT NOT NULL,
    image_url TEXT NOT NULL,
    price DECIMAL NOT NULL,
    name TEXT NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- Comma added here
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (seller_id) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);
