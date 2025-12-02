-- name: CreateRestaurant :one
INSERT INTO "restaurant" (name, description, address, category, city, district, logo_url, banner_url, phone_number, website_url, email, user_id)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id;

-- name: CreateRestaurantHours :exec
INSERT INTO "restaurant_hours" (restaurant_id, day_of_week, open_time, close_time)
VALUES($1, $2, $3, $4);

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
    rs.user_id,
    rsh.day_of_week,
    rsh.open_time,
    rsh.close_time
FROM "restaurant" rs
INNER JOIN "restaurant_hours" rsh ON rs.id = rsh.restaurant_id
WHERE name LIKE $1;

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
    rs.user_id,
    rsh.day_of_week,
    rsh.open_time,
    rsh.close_time
FROM "restaurant" rs
INNER JOIN "restaurant_hours" rsh ON rs.id = rsh.restaurant_id
WHERE id = $1;

-- name: UpdateRestaurant :exec
UPDATE "restaurant"
SET name = $1, description = $2, address = $3,
    category = $4, city = $5, district = $6,
    logo_url = $7, banner_url = $8, phone_number = $9,
    website_url = $10, email = $11, user_id = $12
WHERE id = $13;

-- name: DeleteRestaurantHours :exec
DELETE FROM "restaurant_hours" WHERE restaurant_id = $1;


-- name: DeleteRestaurant :exec
DELETE FROM "restaurant" WHERE id = $1;
