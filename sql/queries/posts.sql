-- name: CreatePost :one
INSERT INTO posts (created_at, updated_at, body, user_id)
VALUES (
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetPosts :many
SELECT
    *
FROM posts
ORDER BY created_at ASC;

-- name: GetPost :one
SELECT
    *
FROM posts
WHERE id = $1
LIMIT 1;

-- name: ResetPosts :exec
TRUNCATE TABLE posts;