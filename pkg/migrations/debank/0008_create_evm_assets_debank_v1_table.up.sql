CREATE TABLE IF NOT EXISTS evm_assets_debank_v1 (
    evm_asset_id SERIAL PRIMARY KEY,
    wallet_id INTEGER NOT NULL,
    total_usd_value FLOAT,
    chain_list_json TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (wallet_id) REFERENCES global_wallets(wallet_id)
);