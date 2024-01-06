CREATE TABLE IF NOT EXISTS token_prices (
    price_id SERIAL PRIMARY KEY,
    token_mint VARCHAR(255) NOT NULL,
    usd_price DECIMAL,
    exchange_name VARCHAR(255),
    exchange_address VARCHAR(255),
    native_price_value BIGINT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (token_mint) REFERENCES tokens(mint)
);