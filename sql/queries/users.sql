-- name: CreateUser :one
INSERT INTO users (email, hashed_password) VALUES ($1, $2) RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUser :one
SELECT * FROM users WHERE users.id = $1;

-- name: GetByEmail :one
SELECT * FROM users WHERE users.email = $1;

-- name: UpdateUser :one
UPDATE users SET email = $2, hashed_password = $3 WHERE id = $1 RETURNING *;

-- name: UpgradeUserRed :one
UPDATE users SET is_chirpy_red = true WHERE id = $1 RETURNING *;