CREATE TABLE IF NOT EXISTS nft_list (
    attribute_id SERIAL PRIMARY KEY,
    wallet_address VARCHAR(255) REFERENCES evm_assets_debank_v1(user_address),
    nft_id VARCHAR(255),
    key VARCHAR(255),
    trait_type VARCHAR(255),
    value TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
