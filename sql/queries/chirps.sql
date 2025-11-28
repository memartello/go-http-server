-- name: CreateChirp :one
INSERT INTO chirp (user_id, body)  VALUES ($1, $2) RETURNING *;

-- name: GetChirp :one
SELECT * FROM chirp WHERE chirp.id = $1 LIMIT 1;

-- name: GetChirps :many
SELECT * FROM chirp ORDER BY chirp.created_at ASC;

-- name: DeleteChirp :exec
DELETE FROM chirp WHERE id = $1;