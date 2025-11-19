-- name: GetUserByID :one
SELECT u.id, u.email, u.full_name, r.role_name, u.is_active, u.created_at, u.updated_at FROM "user" u
INNER JOIN "role" r ON r.role_id = u.role_id
WHERE u.id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT u.id, u.email, u.full_name, r.role_name, u.is_active, u.created_at, u.updated_at FROM "user" u
INNER JOIN  "role" r ON r.role_id = u.role_id
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO "user" (email, full_name, password_hash, role_id) VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: UpdateUser :one
UPDATE "user"
SET full_name = $1, email = $2, password_hash = $3, is_active = $4, updated_at = NOW()
WHERE id = $5
RETURNING *;