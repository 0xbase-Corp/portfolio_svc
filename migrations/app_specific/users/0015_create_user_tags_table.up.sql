CREATE TABLE IF NOT EXISTS user_tags (
    tag_link_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    portfolio_id INTEGER,
    token_id INTEGER,
    nft_id INTEGER,
    tag_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (portfolio_id) REFERENCES portfolios(portfolio_id),
    FOREIGN KEY (token_id) REFERENCES tokens(token_id),
    FOREIGN KEY (nft_id) REFERENCES nfts(nft_id),
    FOREIGN KEY (tag_id) REFERENCES tags(tag_id)
    -- Ensure that portfolios, tokens, nfts, and tags tables are created before this script is run
);
