CREATE TABLE IF NOT EXISTS pseudonymous_portfolios (
    portfolio_id SERIAL PRIMARY KEY,
    cryptographic_link VARCHAR(255) NOT NULL,
    unique_portfolio_identifier VARCHAR(255) NOT NULL,
    portfolio_name VARCHAR(255),
    portfolio_category VARCHAR(255),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);