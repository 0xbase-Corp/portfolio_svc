CREATE TABLE IF NOT EXISTS chain_details (
    chain_id VARCHAR(255) PRIMARY KEY,
    community_id VARCHAR(255),
    name VARCHAR(255),
    logo_url VARCHAR(255),
    native_token_id VARCHAR(255),
    wrapped_token_id VARCHAR(255),
    usd_value DECIMAL,
    wallet_id INTEGER NOT NULL UNIQUE,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (wallet_id) REFERENCES global_wallets(wallet_id) ON DELETE CASCADE
);