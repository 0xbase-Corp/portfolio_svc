CREATE TABLE IF NOT EXISTS tokens (
    token_id SERIAL PRIMARY KEY,
    solana_asset_id INTEGER NOT NULL,
    associated_token_address VARCHAR(255),
    mint VARCHAR(255), 
    amount_raw VARCHAR(255),
    amount VARCHAR(255),
    decimals VARCHAR(255),
    name VARCHAR(255), -- New column for token name
    symbol VARCHAR(50), -- New column for token symbol
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (solana_asset_id) REFERENCES solana_assets_moralis_v1(solana_asset_id)
);