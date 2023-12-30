CREATE TABLE IF NOT EXISTS bitcoin_address_info (
    address VARCHAR(255) PRIMARY KEY,
    received FLOAT,
    sent FLOAT,
    balance FLOAT,
    tx_count INTEGER,
    unconfirmed_tx_count INTEGER,
    unconfirmed_received FLOAT,
    unconfirmed_sent FLOAT,
    unspent_tx_count INTEGER,
    first_tx TEXT,
    last_tx TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);