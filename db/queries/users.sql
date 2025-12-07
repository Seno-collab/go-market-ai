-- name: GetUserByID :one
SELECT u.id, u.email, u.full_name, r.role_name, u.is_active, u.created_at, u.updated_at FROM "users" u
LEFT JOIN "roles" r ON r.id = u.role_id
WHERE u.id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT u.id, u.email, u.full_name, r.role_name, u.password_hash, u.is_active, u.created_at, u.updated_at FROM "users" u
LEFT JOIN  "roles" r ON r.id = u.role_id
WHERE email = $1 LIMIT 1;

-- name: GetUserByName :one
SELECT u.id, u.email, u.full_name, r.role_name, u.is_active, u.created_at, u.updated_at FROM "users" u
LEFT JOIN  "roles" r ON r.id = u.role_id
WHERE full_name = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO "users" (email, full_name, password_hash, role_id) VALUES ($1, $2, $3,(SELECT id FROM "roles" WHERE role_name = 'user'))
RETURNING id;

-- name: UpdateUser :exec
UPDATE "users"
SET full_name = $1, email = $2, password_hash = $3, is_active = $4, updated_at = NOW()
WHERE id = $5;

-- name: GetUserRole :one
SELECT  r.role_name FROM "users" u
LEFT JOIN  "roles" r ON r.id = u.role_id
WHERE u.id = $1 AND u.is_active = $2 LIMIT 1;

-- name: GetPasswordByID :one
SELECT password_hash
FROM "users"
WHERE id = $1;


-- name: UpdatePasswordByID :exec
UPDATE "users"
SET password_hash = $1
WHERE id = $2;
