ALTER TABLE products
ALTER COLUMN stock SET DATA TYPE INTEGER,
ADD CONSTRAINT stock_positive CHECK (stock >= 0);

-- Add an index to the name column
CREATE INDEX idx_products_name ON products (name);