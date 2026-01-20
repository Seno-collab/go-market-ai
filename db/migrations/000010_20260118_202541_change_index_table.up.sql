BEGIN;

-- TOPICS
DROP INDEX IF EXISTS idx_topic_parent;
CREATE INDEX idx_topics_parent_id ON topics(parent_id);

-- MENU_ITEMS
DROP INDEX IF EXISTS idx_menu_item_restaurant;
DROP INDEX IF EXISTS idx_menu_item_topic;
DROP INDEX IF EXISTS idx_menu_item_type;

CREATE INDEX idx_menu_items_restaurant_id ON menu_items(restaurant_id);
CREATE INDEX idx_menu_items_topic_id ON menu_items(topic_id);
CREATE INDEX idx_menu_items_type ON menu_items(type);

-- MENU_ITEM_VARIANTS
DROP INDEX IF EXISTS idx_variant_item;
CREATE INDEX idx_menu_item_variants_menu_item_id
ON menu_item_variants(menu_item_id);

-- OPTION_ITEMS
DROP INDEX IF EXISTS idx_option_group;
CREATE INDEX idx_option_items_option_group_id
ON option_items(option_group_id);

-- COMBO_GROUPS
DROP INDEX IF EXISTS idx_combo_group_combo;
CREATE INDEX idx_combo_groups_combo_item_id
ON combo_groups(combo_item_id);

-- COMBO_GROUP_ITEMS
DROP INDEX IF EXISTS idx_combo_group_item_group;
CREATE INDEX idx_combo_group_items_combo_group_id
ON combo_group_items(combo_group_id);

COMMIT;
