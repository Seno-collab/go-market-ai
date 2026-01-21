-- name: GetTopicsByRestaurant :many
SELECT *
FROM topics
WHERE restaurant_id = $1 AND ( $2::text = '' OR name = $2)
ORDER BY sort_order, id
LIMIT $3 OFFSET $4;


-- name: CountTopicsByRestaurant :one
SELECT COUNT(*) FROM topics WHERE restaurant_id = $1;

-- name: GetTopic :one
SELECT * FROM topics WHERE id = $1 AND restaurant_id = $2;


-- name: GetTopicsByRestaurantCombobox :many
SELECT id as Value, name as TEXT FROM topics WHERE restaurant_id = $1 and is_active = true;

-- name: CreateTopic :one
INSERT INTO topics (restaurant_id, name, slug, parent_id, sort_order)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateTopic :exec
UPDATE topics
SET
    name = $3,
    slug = $4,
    parent_id = $5,
    sort_order = $6,
    updated_at = NOW()
WHERE id = $1
  AND restaurant_id = $2;

-- name: DeleteTopic :exec
DELETE FROM topics
WHERE id = $1
  AND restaurant_id = $2;

-- name: GetMenuItemsByRestaurant :many
SELECT *
FROM menu_items
WHERE restaurant_id = sqlc.arg(restaurant_id)
  AND (
          sqlc.narg(is_active)::boolean IS NULL
          OR is_active = sqlc.narg(is_active)::boolean
      )
  AND (
          sqlc.arg(name)::text = ''
          OR name ILIKE '%' || sqlc.arg(name) || '%'
          OR description ILIKE '%' || sqlc.arg(name) || '%'
        ) AND (NULLIF(sqlc.arg(type)::text, '') IS NULL
        OR type = sqlc.arg(type)::menu_item_type)
ORDER BY sort_order, id
LIMIT sqlc.arg(limit_value) OFFSET sqlc.arg(offset_value);

-- name: CountMenuItems :one
SELECT COUNT(*)
FROM menu_items
WHERE restaurant_id = sqlc.arg(restaurant_id) AND (
    sqlc.narg(is_active)::boolean IS NULL
    OR is_active = sqlc.narg(is_active)::boolean
  )
  AND (
          sqlc.narg(name)::text = ''
          OR name ILIKE '%' || sqlc.arg(name) || '%'
          OR description ILIKE '%' || sqlc.arg(name) || '%'
        ) AND (NULLIF(sqlc.narg(type)::text, '') IS NULL
        OR type = sqlc.narg(type)::menu_item_type);

-- name: GetMenuItemByID :one
SELECT *
FROM menu_items
WHERE id = $1 and restaurant_id = $2;

-- name: CreateMenuItem :one
INSERT INTO menu_items (
    restaurant_id, topic_id, type, name, description,
    image_url, sku, base_price, is_active, sort_order
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, TRUE, $9)
RETURNING id;

-- name: UpdateMenuItem :exec
UPDATE menu_items
SET
    topic_id = $2,
    type = $3,
    name = $4,
    description = $5,
    image_url = $6,
    sku = $7,
    base_price = $8,
    is_active = $9,
    sort_order = $10,
    updated_at = NOW()
WHERE id = $1 and restaurant_id = $11;


-- name: UpdateStatusMenuItem :exec
UPDATE menu_items
SET is_active = $1,
updated_at = NOW()
WHERE id = $2 and restaurant_id = $3;


-- name: DeleteMenuItem :exec
DELETE FROM menu_items WHERE id = $1 AND restaurant_id = $2;

-- name: GetVariantsByItem :many
SELECT *
FROM menu_item_variants
WHERE menu_item_id = $1
ORDER BY sort_order, id;

-- name: CreateVariant :exec
INSERT INTO menu_item_variants (
    menu_item_id, name, price_delta, is_default, sort_order
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UpdateVariant :exec
UPDATE menu_item_variants
SET
    name = $2,
    price_delta = $3,
    is_default = $4,
    sort_order = $5,
    updated_at = NOW()
WHERE id = $1;
DELETE FROM menu_item_variants WHERE id = $1;

-- name: GetOptionGroupsByItem :many
SELECT og.*
FROM option_groups og
JOIN menu_item_option_groups mig ON mig.option_group_id = og.id
WHERE mig.menu_item_id = $1
  AND og.restaurant_id = $2
ORDER BY og.sort_order, og.id
LIMIT $3 OFFSET $4;

-- name: CreateOptionGroup :one
INSERT INTO option_groups (
    restaurant_id, name, min_select, max_select, is_required, sort_order
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: AttachOptionGroupToItem :exec
INSERT INTO menu_item_option_groups(menu_item_id, option_group_id, sort_order)
VALUES ($1, $2, $3)
ON CONFLICT (menu_item_id, option_group_id) DO NOTHING;

-- name: GetOptionGroup :one
SELECT *
FROM option_groups
WHERE id = $1
  AND restaurant_id = $2;

-- name: UpdateOptionGroup :exec
UPDATE option_groups
SET
    name = $3,
    min_select = $4,
    max_select = $5,
    is_required = $6,
    sort_order = $7,
    updated_at = NOW()
WHERE id = $1
  AND restaurant_id = $2;

-- name: DeleteOptionGroup :exec
DELETE FROM option_groups
WHERE id = $1
  AND restaurant_id = $2;

-- name: GetOptionItemsByGroup :many
SELECT oi.*
FROM option_items oi
JOIN option_groups og ON og.id = oi.option_group_id
WHERE oi.option_group_id = $1
  AND og.restaurant_id = $2
ORDER BY oi.sort_order, oi.id
LIMIT $3 OFFSET $4;

-- name: CountOptionItems :one
SELECT COUNT(*)
FROM option_items oi
JOIN option_groups og ON og.id = oi.option_group_id
WHERE oi.option_group_id = $1
  AND og.restaurant_id = $2
ORDER BY oi.sort_order, oi.id;

-- name: GetOptionItem :one
SELECT oi.*
FROM option_items oi
JOIN option_groups og ON og.id = oi.option_group_id
WHERE oi.id = $1
  AND og.restaurant_id = $2;

-- name: CreateOptionItem :one
INSERT INTO option_items (
    option_group_id, name, linked_menu_item,
    price_delta, quantity_min, quantity_max, sort_order
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: UpdateOptionItem :exec
UPDATE option_items oi
SET
    name = $2,
    linked_menu_item = $3,
    price_delta = $4,
    quantity_min = $5,
    quantity_max = $6,
    sort_order = $7,
    updated_at = NOW()
FROM option_groups og
WHERE oi.id = $1
  AND oi.option_group_id = og.id
  AND og.restaurant_id = $8;

-- name: DeleteOptionItem :exec
DELETE FROM option_items oi
USING option_groups og
WHERE oi.id = $1
  AND oi.option_group_id = og.id
  AND og.restaurant_id = $2;

-- name: CreateComboGroup :exec
INSERT INTO combo_groups (
    combo_item_id, name, min_select, max_select, sort_order
)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateComboGroup :exec
UPDATE combo_groups
SET
    name = $2,
    min_select = $3,
    max_select = $4,
    sort_order = $5,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteComboGroup :exec
DELETE FROM combo_groups WHERE id = $1;

-- name: GetComboGroupItems :many
SELECT *
FROM combo_group_items
WHERE combo_group_id = $1
ORDER BY sort_order, id;

-- name: CreateComboGroupItem :exec
INSERT INTO combo_group_items (
    combo_group_id, menu_item_id,
    price_delta, quantity_default, quantity_min,
    quantity_max, sort_order
)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: DeleteComboGroupItem :exec
DELETE FROM combo_group_items WHERE id = $1;


-- name: ListMenus :many
SELECT
    m.*
FROM menu_items m
JOIN topics t ON m.topic_id = t.id
WHERE m.restaurant_id = sqlc.arg(restaurant_id)
  AND (NULLIF(sqlc.arg(menu_type)::text, '') IS NULL OR m.type = sqlc.arg(menu_type)::menu_item_type)
  AND (NULLIF(sqlc.arg(cursor)::bigint, 0) IS NULL OR m.id < sqlc.arg(cursor))
  AND (COALESCE(sqlc.arg(topic_names)::text[], '{}'::text[]) = '{}' OR t.name = ANY(sqlc.arg(topic_names)))
  AND m.is_active = true
GROUP BY m.id, m.name, m.restaurant_id, m.type
HAVING (COALESCE(sqlc.arg(topic_names)::text[], '{}'::text[]) = '{}' OR COUNT(DISTINCT t.name) = cardinality(sqlc.arg(topic_names)))
ORDER BY m.id DESC
LIMIT sqlc.arg(page_size);


-- name: GetFullMenu :many
SELECT
    mi.*,
    t.id AS topic_id,
    t.name AS topic_name,
    mv.*,
    og.*,
    oi.*,
    cg.*,
    cgi.*
FROM menu_items mi
LEFT JOIN topics t ON t.id = mi.topic_id
LEFT JOIN menu_item_variants mv ON mv.menu_item_id = mi.id
LEFT JOIN menu_item_option_groups mog ON mog.menu_item_id = mi.id
LEFT JOIN option_groups og ON og.id = mog.option_group_id
LEFT JOIN option_items oi ON oi.option_group_id = og.id
LEFT JOIN combo_groups cg ON cg.combo_item_id = mi.id
LEFT JOIN combo_group_items cgi ON cgi.combo_group_id = cg.id
WHERE mi.restaurant_id = $1;
