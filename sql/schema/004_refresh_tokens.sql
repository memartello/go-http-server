-- +goose Up

CREATE TABLE IF NOT EXISTS refresh_tokens (
    token TEXT PRIMARY KEY NOT NULL,
    user_id uuid NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
	revoked_at TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down

DROP TABLE IF EXISTS refresh_tokens;