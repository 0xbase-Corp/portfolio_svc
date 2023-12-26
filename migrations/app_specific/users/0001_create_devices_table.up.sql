CREATE TABLE IF NOT EXISTS devices (
    device_id VARCHAR(255) PRIMARY KEY,
    user_id INTEGER NOT NULL,
    device_type VARCHAR(255),
    device_identifier VARCHAR(255),
    persistent_login_token TEXT,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
    -- Make sure the users table is created with a user_id column before running this script.
);
