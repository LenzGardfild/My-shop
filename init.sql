-- Создание таблиц для хранения заказов
CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(64) PRIMARY KEY,
    track_number VARCHAR(64),
    entry VARCHAR(16),
    locale VARCHAR(8),
    internal_signature VARCHAR(64),
    customer_id VARCHAR(64),
    delivery_service VARCHAR(64),
    shardkey VARCHAR(8),
    sm_id INTEGER,
    date_created TIMESTAMP,
    oof_shard VARCHAR(8)
);

CREATE TABLE IF NOT EXISTS deliveries (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(64) REFERENCES orders(order_uid) ON DELETE CASCADE,
    name VARCHAR(128),
    phone VARCHAR(32),
    zip VARCHAR(32),
    city VARCHAR(64),
    address VARCHAR(128),
    region VARCHAR(64),
    email VARCHAR(128)
);

CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(64) REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction VARCHAR(64),
    request_id VARCHAR(64),
    currency VARCHAR(8),
    provider VARCHAR(32),
    amount INTEGER,
    payment_dt BIGINT,
    bank VARCHAR(64),
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER
);

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(64) REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id BIGINT,
    track_number VARCHAR(64),
    price INTEGER,
    rid VARCHAR(64),
    name VARCHAR(128),
    sale INTEGER,
    size VARCHAR(16),
    total_price INTEGER,
    nm_id BIGINT,
    brand VARCHAR(64),
    status INTEGER
);
