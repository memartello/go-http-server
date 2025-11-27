-- name: CreateUser :one
INSERT INTO users (email) VALUES ($1) RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUser :one
SELECT * FROM users WHERE users.id = $1;