CREATE TABLE IF NOT EXISTS global_wallets (
    wallet_id SERIAL PRIMARY KEY,
    portfolio_id INTEGER NOT NULL,
    blockchain_type VARCHAR(255) NOT NULL,
    api_endpoint TEXT,
    api_version VARCHAR(50),
    last_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (portfolio_id) REFERENCES pseudonymous_portfolios(portfolio_id)
);
