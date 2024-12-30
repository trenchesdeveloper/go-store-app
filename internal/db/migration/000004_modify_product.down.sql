ALTER TABLE products
DROP CONSTRAINT IF EXISTS stock_positive;

-- Remove the index on the name column
DROP INDEX IF EXISTS idx_products_name;