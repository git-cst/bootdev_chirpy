-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email)
VALUES (
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: GetUser :one
SELECT 
    id,
    created_at,
    updated_at,
    email
FROM users
WHERE id = $1
LIMIT 1;

-- name: ResetUsers :exec
TRUNCATE TABLE users CASCADE;