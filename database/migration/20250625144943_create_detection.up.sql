CREATE TABLE detections (
    id SERIAL NOT NULL PRIMARY KEY, 
    ip_address INET NOT NULL,
    customer_id UUID NOT NULL,
    path VARCHAR(128) NOT NULL DEFAULT '',
    method VARCHAR(16) NOT NULL DEFAULT '',
    status_code SMALLINT NOT NULL,
    type VARCHAR(32) NOT NULL default 'unknown',
    confidence NUMERIC(10,9),
    is_banned BOOLEAN NOT NULL default false,
    is_deepfake BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,

    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

CREATE INDEX idx_detections_customer_id ON detections(customer_id);

CREATE INDEX idx_detections_customer_time ON detections(customer_id, created_at DESC);
