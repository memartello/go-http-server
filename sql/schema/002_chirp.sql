-- +goose Up

CREATE TABLE IF NOT EXISTS chirp (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id uuid NOT NULL,
	body TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS chirp;