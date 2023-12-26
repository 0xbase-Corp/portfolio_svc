
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_public_key TEXT,
    other_user_details TEXT,
    signup_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

