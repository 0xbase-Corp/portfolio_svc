CREATE TABLE IF NOT EXISTS nft_list (
    attribute_id SERIAL PRIMARY KEY,
    evm_asset_id INTEGER NOT NULL,
    nft_id VARCHAR(255),
    key VARCHAR(255),
    trait_type VARCHAR(255),
    value TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (evm_asset_id) REFERENCES evm_assets_debank_v1(evm_asset_id)
);