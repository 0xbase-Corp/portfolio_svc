CREATE TABLE IF NOT EXISTS chain_details (
    chain_id VARCHAR(255) PRIMARY KEY,
    community_id VARCHAR(255),
    name VARCHAR(255),
    logo_url VARCHAR(255),
    native_token_id VARCHAR(255),
    wrapped_token_id VARCHAR(255),
    usd_value DECIMAL,
    user_address VARCHAR(255) REFERENCES evm_assets_debank_v1(user_address),
    updated_at TIMESTAMP,
    created_at TIMESTAMP
);
