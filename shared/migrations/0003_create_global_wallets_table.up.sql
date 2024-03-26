CREATE TABLE IF NOT EXISTS global_wallets (
    wallet_id SERIAL PRIMARY KEY,
    portfolio_id INTEGER NOT NULL,
    wallet_address VARCHAR(255) unique NOT NULL,
    blockchain_type VARCHAR(50) NOT NULL,
    api_endpoint TEXT,
    api_version VARCHAR(50),
    last_updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
