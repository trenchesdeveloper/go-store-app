CREATE TABLE categories (
                            id SERIAL PRIMARY KEY,
                            name VARCHAR(255) NOT NULL,
                            parent_id INTEGER,
                            image_url VARCHAR(255),
                            display_order INTEGER,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(255) NOT NULL,
                          description TEXT,
                          category_id INTEGER NOT NULL,
                          image_url VARCHAR(255),
                          price NUMERIC(10, 2) NOT NULL,
                          user_id INTEGER NOT NULL,
                          stock INTEGER NOT NULL,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

ALTER TABLE products ADD FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE SET NULL;