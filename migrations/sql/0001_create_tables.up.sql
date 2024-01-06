CREATE TABLE IF NOT EXISTS user_profile(
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    google_user_id VARCHAR(100) NOT NULL,
    UNIQUE(google_user_id)
);

CREATE TABLE IF NOT EXISTS user_google_token(
    user_id VARCHAR(255) PRIMARY KEY references user_profile(id),
    access_token VARCHAR(300) NOT NULL,
    refresh_token VARCHAR(300) NOT NULL,
    expiry timestamptz
);