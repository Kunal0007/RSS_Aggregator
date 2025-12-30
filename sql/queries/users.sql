-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, api_key)
VALUES ($1, $2, $3, $4,
    encode(sha256(random()::text::bytea), 'hex')
)
RETURNING *;

-- -- name: UpdateUser :one
-- UPDATE users
-- SET updated_at = $1, name = $2
-- WHERE id = $3
-- RETURNING *;