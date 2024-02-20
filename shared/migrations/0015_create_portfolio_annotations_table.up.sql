CREATE TABLE IF NOT EXISTS portfolio_annotations (
    annotation_id SERIAL PRIMARY KEY,
    portfolio_id INTEGER NOT NULL,
    content TEXT,
    tag TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  
    FOREIGN KEY (portfolio_id) REFERENCES pseudonymous_portfolios(portfolio_id)
);