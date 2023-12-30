CREATE TABLE IF NOT EXISTS user_tags (
    tag_link_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    portfolio_id INTEGER,
    token_id INTEGER,
    nft_id INTEGER,
    tag TEXT,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (portfolio_id) REFERENCES pseudonymous_portfolios(portfolio_id),
    FOREIGN KEY (token_id) REFERENCES tokens(token_id),
    FOREIGN KEY (nft_id) REFERENCES nfts(nft_id),
    CHECK (
        (portfolio_id IS NOT NULL)::integer +
        (token_id IS NOT NULL)::integer +
        (nft_id IS NOT NULL)::integer = 1
    )
);