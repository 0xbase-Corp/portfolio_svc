CREATE TABLE IF NOT EXISTS evm_assets_debank_v1 (
    user_address VARCHAR(255) PRIMARY KEY,
    total_usd_value FLOAT,
    chain_list_json TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
