CREATE TABLE IF NOT EXISTS bitcoin_btc_com_v1 (
    btc_asset_id SERIAL PRIMARY KEY,
    wallet_id INTEGER NOT NULL,
    btc_usd_price FLOAT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
