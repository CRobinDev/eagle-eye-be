CREATE TABLE payments (
    user_id UUID NOT NULL,
    order_id VARCHAR(36) NOT NULL PRIMARY KEY,
    amount DECIMAL(10, 2) NOT NULL,
    tier VARCHAR(32) NOT NULL DEFAULT 'free' CHECK (tier IN ('free', 'basic', 'premium', 'enterprise')),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_payments_deleted_at ON payments (deleted_at);