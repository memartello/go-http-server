-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetUserByToken :one
SELECT refresh_tokens.*, users.* FROM refresh_tokens JOIN users ON users.id = refresh_tokens.user_id WHERE refresh_tokens.token = $1 LIMIT 1;

-- name: RevokeToken :exec
UPDATE refresh_tokens SET revoked_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE token = $1;