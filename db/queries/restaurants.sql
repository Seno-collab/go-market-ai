-- name: CreateRestaurant :one
INSERT INTO "restaurants" (name, description, address, category, city, district, logo_url, banner_url, phone_number, website_url, email, created_by)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id;

-- name: CreateRestaurantHours :exec
INSERT INTO "restaurant_hours" (restaurant_id, day_of_week, open_time, close_time)
VALUES($1, $2, $3, $4);

-- name: UpsertRestaurantUser :exec
INSERT INTO restaurant_users (
    restaurant_id, user_id, role, created_at, updated_at
)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT (restaurant_id, user_id)
DO UPDATE SET
    role = EXCLUDED.role,
    updated_at = NOW();

-- name: CheckUserRestaurant :one
SELECT 1
FROM restaurant_users
WHERE restaurant_id = $1
  AND user_id = $2
  AND deleted_at IS NULL
LIMIT 1;

-- name: SoftDeleteRestaurantUser :exec
UPDATE restaurant_users
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE restaurant_id = $1 AND user_id = $2;

-- name: GetByName :many
SELECT
    rs.id,
    rs.name,
    rs.description,
    rs.address,
    rs.category,
    rs.city,
    rs.district,
    rs.logo_url,
    rs.banner_url,
    rs.phone_number,
    rs.website_url,
    rs.email,
    rs.created_by,
    rsh.day_of_week,
    rsh.open_time,
    rsh.close_time
FROM "restaurants" rs
INNER JOIN "restaurant_hours" rsh ON rs.id = rsh.restaurant_id
WHERE name LIKE $1 AND rs.deleted_at IS NULL;

-- name: GetById :many
SELECT
    rs.id,
    rs.name,
    rs.description,
    rs.address,
    rs.category,
    rs.city,
    rs.district,
    rs.logo_url,
    rs.banner_url,
    rs.phone_number,
    rs.website_url,
    rs.email,
    rs.created_by,
    rsh.day_of_week,
    rsh.open_time,
    rsh.close_time
FROM "restaurants" rs
INNER JOIN "restaurant_hours" rsh ON rs.id = rsh.restaurant_id
WHERE id = $1 AND rs.deleted_at IS NULL;

-- name: UpdateRestaurant :exec
UPDATE "restaurants"
SET name = $1, description = $2, address = $3,
    category = $4, city = $5, district = $6,
    logo_url = $7, banner_url = $8, phone_number = $9,
    website_url = $10, email = $11, updated_by = $12
WHERE id = $13;

-- name: SoftDeleteRestaurantHours :exec
UPDATE "restaurant_hours"
SET deleted_at = NOW()
WHERE restaurant_id = $1;

-- name: SoftDeleteRestaurant :exec
UPDATE restaurants
SET deleted_at = NOW(), updated_by = $1
WHERE id = $2;


-- name: GetRestaurantByUserID :one
SELECT restaurant_id FROM restaurant_users WHERE user_id = $1;