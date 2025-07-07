BEGIN;

-- Delete only seeded detections
DELETE FROM detections
WHERE email IN ('alice@example.com', 'bob@example.com', 'carol@example.com', 'dave@example.com', 'eve@example.com');

-- Delete seeded customers (based on hashed_key or email mapping)
DELETE FROM customers
WHERE hashed_key IN ('key001', 'key002', 'key003', 'key004', 'key005');

-- Delete seeded payments
DELETE FROM payments
WHERE order_id IN ('order-001', 'order-002', 'order-003', 'order-004', 'order-005');

-- Delete seeded users
DELETE FROM users
WHERE email IN ('alice@example.com', 'bob@example.com', 'carol@example.com', 'dave@example.com', 'eve@example.com');

-- Drop the ULID generation function
DROP FUNCTION IF EXISTS generate_ulid_at_time(TIMESTAMP WITH TIME ZONE);

COMMIT;