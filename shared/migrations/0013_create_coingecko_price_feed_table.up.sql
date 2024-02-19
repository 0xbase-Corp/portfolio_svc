CREATE TABLE IF NOT EXISTS coingecko_price_feed (
    crypto_id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    ticker VARCHAR(255),
    usd_value FLOAT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
