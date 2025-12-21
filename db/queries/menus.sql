-- name: GetTopicsByRestaurant :many
SELECT *
FROM topic
WHERE restaurant_id = $1
ORDER BY sort_order, id;


-- name: GetTopic :one
SELECT * FROM topic WHERE id = $1 AND restaurant_id = $2;


-- name: CreateTopic :one
INSERT INTO topic (restaurant_id, name, slug, parent_id, sort_order)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateTopic :exec
UPDATE topic
SET
    name = $3,
    slug = $4,
    parent_id = $5,
    sort_order = $6,
    updated_at = NOW()
WHERE id = $1
  AND restaurant_id = $2;

-- name: DeleteTopic :exec
DELETE FROM topic
WHERE id = $1
  AND restaurant_id = $2;

-- name: GetMenuItemsByRestaurant :many
SELECT *
FROM menu_item
WHERE restaurant_id = $1
ORDER BY sort_order, id;

-- name: GetMenuItemByID :one
SELECT *
FROM menu_item
WHERE id = $1 and restaurant_id = $2;

-- name: CreateMenuItem :one
INSERT INTO menu_item (
    restaurant_id, topic_id, type, name, description,
    image_url, sku, base_price, is_active, sort_order
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, TRUE, $9)
RETURNING id;

-- name: UpdateMenuItem :exec
UPDATE menu_item
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

-- name: DeleteMenuItem :exec
DELETE FROM menu_item WHERE id = $1 AND restaurant_id = $2;

-- name: GetVariantsByItem :many
SELECT *
FROM menu_item_variant
WHERE menu_item_id = $1
ORDER BY sort_order, id;

-- name: CreateVariant :exec
INSERT INTO menu_item_variant (
    menu_item_id, name, price_delta, is_default, sort_order
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UpdateVariant :exec
UPDATE menu_item_variant
SET
    name = $2,
    price_delta = $3,
    is_default = $4,
    sort_order = $5,
    updated_at = NOW()
WHERE id = $1;
DELETE FROM menu_item_variant WHERE id = $1;

-- name: GetOptionGroupsByItem :many
SELECT og.*
FROM option_group og
JOIN menu_item_option_group mig ON mig.option_group_id = og.id
WHERE mig.menu_item_id = $1
  AND og.restaurant_id = $2
ORDER BY og.sort_order, og.id;

-- name: CreateOptionGroup :one
INSERT INTO option_group (
    restaurant_id, name, min_select, max_select, is_required, sort_order
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: AttachOptionGroupToItem :exec
INSERT INTO menu_item_option_group (menu_item_id, option_group_id, sort_order)
VALUES ($1, $2, $3)
ON CONFLICT (menu_item_id, option_group_id) DO NOTHING;

-- name: GetOptionGroup :one
SELECT *
FROM option_group
WHERE id = $1
  AND restaurant_id = $2;

-- name: UpdateOptionGroup :exec
UPDATE option_group
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
DELETE FROM option_group
WHERE id = $1
  AND restaurant_id = $2;

-- name: GetOptionItemsByGroup :many
SELECT oi.*
FROM option_item oi
JOIN option_group og ON og.id = oi.option_group_id
WHERE oi.option_group_id = $1
  AND og.restaurant_id = $2
ORDER BY oi.sort_order, oi.id;

-- name: GetOptionItem :one
SELECT oi.*
FROM option_item oi
JOIN option_group og ON og.id = oi.option_group_id
WHERE oi.id = $1
  AND og.restaurant_id = $2;

-- name: CreateOptionItem :one
INSERT INTO option_item (
    option_group_id, name, linked_menu_item,
    price_delta, quantity_min, quantity_max, sort_order
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: UpdateOptionItem :exec
UPDATE option_item oi
SET
    name = $2,
    linked_menu_item = $3,
    price_delta = $4,
    quantity_min = $5,
    quantity_max = $6,
    sort_order = $7,
    updated_at = NOW()
FROM option_group og
WHERE oi.id = $1
  AND oi.option_group_id = og.id
  AND og.restaurant_id = $8;

-- name: DeleteOptionItem :exec
DELETE FROM option_item oi
USING option_group og
WHERE oi.id = $1
  AND oi.option_group_id = og.id
  AND og.restaurant_id = $2;

-- name: CreateComboGroup :exec
INSERT INTO combo_group (
    combo_item_id, name, min_select, max_select, sort_order
)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateComboGroup :exec
UPDATE combo_group
SET
    name = $2,
    min_select = $3,
    max_select = $4,
    sort_order = $5,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteComboGroup :exec
DELETE FROM combo_group WHERE id = $1;

-- name: GetComboGroupItems :many
SELECT *
FROM combo_group_item
WHERE combo_group_id = $1
ORDER BY sort_order, id;

-- name: CreateComboGroupItem :exec
INSERT INTO combo_group_item (
    combo_group_id, menu_item_id,
    price_delta, quantity_default, quantity_min,
    quantity_max, sort_order
)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: DeleteComboGroupItem :exec
DELETE FROM combo_group_item WHERE id = $1;


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
FROM menu_item mi
LEFT JOIN topic t ON t.id = mi.topic_id
LEFT JOIN menu_item_variant mv ON mv.menu_item_id = mi.id
LEFT JOIN menu_item_option_group mog ON mog.menu_item_id = mi.id
LEFT JOIN option_group og ON og.id = mog.option_group_id
LEFT JOIN option_item oi ON oi.option_group_id = og.id
LEFT JOIN combo_group cg ON cg.combo_item_id = mi.id
LEFT JOIN combo_group_item cgi ON cgi.combo_group_id = cg.id
WHERE mi.restaurant_id = $1;
