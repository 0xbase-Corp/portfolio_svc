CREATE TABLE IF NOT EXISTS bitcoin_btc_com_v1 (
    btc_asset_id SERIAL PRIMARY KEY,
    wallet_id INTEGER NOT NULL UNIQUE,
    btc_usd_price FLOAT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (wallet_id) REFERENCES global_wallets(wallet_id) ON DELETE CASCADE
);
