-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email)
VALUES (
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: ResetUsers :exec
TRUNCATE TABLE users CASCADE;