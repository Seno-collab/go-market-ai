-- name: GetUserByID :one
SELECT id, email, name FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, email, name FROM users WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (email, name) VALUES ($1, $2)
RETURNING id, email, name;