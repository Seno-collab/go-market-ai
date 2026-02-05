-- name: GetUserByID :one
SELECT u.id, u.email, u.full_name, r.role_name, u.is_active, u.created_at, u.updated_at, u.image_url
FROM "users" u
LEFT JOIN "roles" r ON r.id = u.role_id
WHERE u.id = sqlc.arg(user_id)::UUID
LIMIT 1;

-- name: GetUserByEmail :one
SELECT u.id, u.email, u.full_name, r.role_name, u.password_hash, u.is_active, u.created_at, u.updated_at, u.image_url
FROM "users" u
LEFT JOIN "roles" r ON r.id = u.role_id
WHERE u.email = sqlc.arg(email)::TEXT
LIMIT 1;

-- name: CreateUser :one
INSERT INTO "users" (email, full_name, password_hash, role_id)
VALUES (
    sqlc.arg(email)::TEXT,
    sqlc.arg(full_name)::TEXT,
    sqlc.arg(password_hash)::TEXT,
    (SELECT id FROM "roles" WHERE role_name = 'user')
)
RETURNING id;

-- name: UpdateUser :exec
UPDATE "users"
SET full_name = sqlc.arg(full_name)::TEXT,
    email = sqlc.arg(email)::TEXT,
    password_hash = sqlc.arg(password_hash)::TEXT,
    image_url = sqlc.arg(image_url)::TEXT,
    is_active = sqlc.arg(is_active)::BOOLEAN,
    updated_at = NOW()
WHERE id = sqlc.arg(user_id)::UUID;

-- name: GetUserRole :one
SELECT r.role_name
FROM "users" u
LEFT JOIN "roles" r ON r.id = u.role_id
WHERE u.id = sqlc.arg(user_id)::UUID
  AND u.is_active = sqlc.arg(is_active)::BOOLEAN
LIMIT 1;

-- name: GetPasswordByID :one
SELECT password_hash
FROM "users"
WHERE id = sqlc.arg(user_id)::UUID;

-- name: UpdatePasswordByID :exec
UPDATE "users"
SET password_hash = sqlc.arg(password_hash)::TEXT
WHERE id = sqlc.arg(user_id)::UUID;
