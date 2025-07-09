-- Ensure pgcrypto is available for gen_random_bytes
DROP EXTENSION IF EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS pgcrypto SCHEMA public;

-- Create a ULID-like UUID generator function
CREATE OR REPLACE FUNCTION generate_ulid_at_time(target_time TIMESTAMP WITH TIME ZONE)
    RETURNS UUID AS
$$
DECLARE
    timestamp_ms BIGINT;
    rand_bytes   BYTEA;
    result       UUID;
BEGIN
    timestamp_ms := FLOOR(EXTRACT(EPOCH FROM target_time) * 1000);
    rand_bytes := gen_random_bytes(10);
    result := CONCAT_WS('-',
                        LPAD(TO_HEX(timestamp_ms >> 16), 8, '0'),
                        LPAD(TO_HEX((timestamp_ms & x'FFFF'::int)), 4, '0'),
                        LPAD(TO_HEX(x'4000'::int | (get_byte(rand_bytes, 0) & x'0FFF'::int)), 4, '0'),
                        LPAD(TO_HEX(x'8000'::int | (get_byte(rand_bytes, 1) & x'3FFF'::int)), 4, '0'),
                        CONCAT(
                            LPAD(TO_HEX(get_byte(rand_bytes, 2)), 2, '0'),
                            LPAD(TO_HEX(get_byte(rand_bytes, 3)), 2, '0'),
                            LPAD(TO_HEX(get_byte(rand_bytes, 4)), 2, '0'),
                            LPAD(TO_HEX(get_byte(rand_bytes, 5)), 2, '0'),
                            LPAD(TO_HEX(get_byte(rand_bytes, 6)), 2, '0'),
                            LPAD(TO_HEX(get_byte(rand_bytes, 7)), 2, '0')
                        )
              )::UUID;
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Now generate dynamic data
DO
$$
DECLARE
    u1 UUID; u2 UUID; u3 UUID; u4 UUID; u5 UUID;
    o1 TEXT := 'order-001'; o2 TEXT := 'order-002'; o3 TEXT := 'order-003'; o4 TEXT := 'order-004'; o5 TEXT := 'order-005';
BEGIN
    -- Generate user IDs with different timestamps
    u1 := generate_ulid_at_time(NOW() - INTERVAL '5 days');
    u2 := generate_ulid_at_time(NOW() - INTERVAL '4 days');
    u3 := generate_ulid_at_time(NOW() - INTERVAL '3 days');
    u4 := generate_ulid_at_time(NOW() - INTERVAL '2 days');
    u5 := generate_ulid_at_time(NOW() - INTERVAL '1 day');

    -- Insert users
    INSERT INTO users (id, username, email, password_hash, role, is_verified, is_customer, auth_provider)
    VALUES
        (u1, 'alice', 'alice@example.com', 'hashed123', 'user', TRUE, TRUE, 'basic'),
        (u2, 'bob', 'bob@example.com', 'hashed123', 'user', TRUE, FALSE, 'google'),
        (u3, 'carol', 'carol@example.com', 'hashed123', 'admin', TRUE, TRUE, 'basic'),
        (u4, 'dave', 'dave@example.com', 'hashed123', 'user', FALSE, TRUE, 'basic'),
        (u5, 'eve', 'eve@example.com', 'hashed123', 'user', TRUE, TRUE, 'google');

    -- Insert payments
    INSERT INTO payments (user_id, order_id, amount, tier, status)
    VALUES
        (u1, o1, 9.99, 'basic', 'completed'),
        (u2, o2, 0.00, 'free', 'pending'),
        (u3, o3, 19.99, 'premium', 'completed'),
        (u4, o4, 49.99, 'enterprise', 'completed'),
        (u5, o5, 9.99, 'basic', 'failed');

    -- Insert customers
    INSERT INTO customers (id, order_id, customer_tier, hashed_key, prefix, current_usage, monthly_limit, expires_at, last_used)
    VALUES
        (u1, o1, 1, 'key001', 'client', 10, 100, NOW() + INTERVAL '30 days', NOW()),
        (u2, o2, 0, 'key002', 'client', 5, 50, NOW() + INTERVAL '30 days', NOW()),
        (u3, o3, 2, 'key003', 'client', 20, 200, NOW() + INTERVAL '60 days', NOW()),
        (u4, o4, 2, 'key004', 'client', 15, 150, NOW() + INTERVAL '30 days', NOW()),
        (u5, o5, 1, 'key005', 'client', 0, 100, NOW() + INTERVAL '30 days', NOW());

    -- Insert detections
    INSERT INTO detections (email, ip_address, customer_id, path, method, status_code, is_deepfake)
    VALUES
        ('alice@example.com', '192.168.1.1', u1, '/api/detect', 'POST', 200, FALSE),
        ('bob@example.com', '192.168.1.2', u2, '/api/detect', 'POST', 500, TRUE),
        ('carol@example.com', '192.168.1.3', u3, '/api/detect', 'GET', 403, FALSE),
        ('dave@example.com', '192.168.1.4', u4, '/api/check', 'POST', 200, TRUE),
        ('eve@example.com', '192.168.1.5', u5, '/api/check', 'POST', 200, FALSE);
END
$$;