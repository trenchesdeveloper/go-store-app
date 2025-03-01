-- Drop the index if it exists
DROP INDEX IF EXISTS idx_users_email;

-- Drop the custom type if it exists
DROP TYPE IF EXISTS user_type CASCADE;

CREATE TYPE user_type AS ENUM ('buyer', 'seller', 'admin');

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       first_name VARCHAR(50) NOT NULL,
                       last_name VARCHAR(50) NOT NULL,
                       email VARCHAR(100) NOT NULL UNIQUE,
                       password VARCHAR(255) NOT NULL,
                       phone VARCHAR(20),
                       code VARCHAR(50),
                       expiry TIMESTAMP,
                       verified BOOLEAN NOT NULL DEFAULT FALSE,
                       user_type user_type NOT NULL DEFAULT 'buyer',
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);