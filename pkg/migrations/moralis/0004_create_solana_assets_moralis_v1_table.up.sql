CREATE TABLE IF NOT EXISTS solana_assets_moralis_v1 (
    solana_asset_id SERIAL PRIMARY KEY,
<<<<<<< HEAD
    wallet_id INTEGER NOT NULL UNIQUE,
=======
    wallet_id INTEGER NOT NULL,
>>>>>>> main
    lamports VARCHAR(255),
    solana varchar(255),
    total_tokens_count INTEGER,
    total_nfts_count INTEGER,
    last_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (wallet_id) REFERENCES global_wallets(wallet_id)
);