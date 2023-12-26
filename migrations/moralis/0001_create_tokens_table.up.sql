CREATE TABLE IF NOT EXISTS tokens (
    token_id SERIAL PRIMARY KEY,
    solana_asset_id INTEGER NOT NULL,
    associated_token_address VARCHAR(255),
    mint VARCHAR(255),
    amount_raw BIGINT,
    amount DECIMAL,
    decimals INTEGER,
    FOREIGN KEY (solana_asset_id) REFERENCES solana_assets_moralis_v1(solana_asset_id)
);
