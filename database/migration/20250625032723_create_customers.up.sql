CREATE TABLE customers (
    id UUID PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL DEFAULT 'free',
    customer_tier SMALLINT NOT NULL CHECK (customer_tier IN (0, 1, 2)), 
    hashed_key VARCHAR(96) NOT NULL,
    prefix VARCHAR(32) NOT NULL DEFAULT 'client',
    current_usage INT NOT NULL DEFAULT 0 CHECK (current_usage <= monthly_limit AND current_usage >= 0),
    monthly_limit INT NOT NULL DEFAULT 100,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMPTZ,
    last_used TIMESTAMPTZ,

    CONSTRAINT fk_customer_payment FOREIGN KEY (order_id) REFERENCES payments(order_id) ON DELETE CASCADE,
    CONSTRAINT fk_customer_user FOREIGN KEY (id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_customers_hashed_key ON customers (hashed_key);
CREATE INDEX IF NOT EXISTS idx_customers_prefix ON customers (prefix);
CREATE INDEX IF NOT EXISTS idx_customers_revoked_at ON customers (revoked_at);
