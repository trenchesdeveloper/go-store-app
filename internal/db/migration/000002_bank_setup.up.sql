CREATE TABLE bank_account (
                       id SERIAL PRIMARY KEY,
                       user_id bigint NOT NULL,
                       bank_account bigint NOT NULL UNIQUE,
                       swift_code VARCHAR(255),
                       payment_type VARCHAR(255),
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE  bank_account ADD FOREIGN KEY (user_id) REFERENCES "users" ("id");

-- Indexes
CREATE INDEX ON "bank_account" ("user_id");