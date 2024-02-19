CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_public_key TEXT,
    other_user_details TEXT,
    signup_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE public.users
  -- taken from OSWAP https://owasp.org/www-community/OWASP_Validation_Regex_Repository
  ADD CONSTRAINT users_email_check CHECK (email ~* '^[a-zA-Z0-9_+&*-]+(?:\.[a-zA-Z0-9_+&*-]+)*@(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$');