-- Drop all restaurant and menu related tables
-- Keep only user and role tables

-- Drop tables in correct order to avoid FK constraint issues

-- Drop combo tables first (both singular and plural)
DROP TABLE IF EXISTS combo_group_item CASCADE;
DROP TABLE IF EXISTS combo_group_items CASCADE;
DROP TABLE IF EXISTS combo_group CASCADE;
DROP TABLE IF EXISTS combo_groups CASCADE;

-- Drop option tables (both singular and plural)
DROP TABLE IF EXISTS option_item CASCADE;
DROP TABLE IF EXISTS option_items CASCADE;
DROP TABLE IF EXISTS menu_item_option_group CASCADE;
DROP TABLE IF EXISTS menu_item_option CASCADE;
DROP TABLE IF EXISTS menu_item_options CASCADE;
DROP TABLE IF EXISTS option_group CASCADE;
DROP TABLE IF EXISTS option_groups CASCADE;

-- Drop menu item variant (both singular and plural)
DROP TABLE IF EXISTS menu_item_variant CASCADE;
DROP TABLE IF EXISTS menu_item_variants CASCADE;

-- Drop menu item (both singular and plural naming)
DROP TABLE IF EXISTS menu_item CASCADE;
DROP TABLE IF EXISTS menu_items CASCADE;

-- Drop enum type for menu_item (CASCADE to drop dependent objects)
DROP TYPE IF EXISTS menu_item_type CASCADE;

-- Drop topic (both singular and plural)
DROP TABLE IF EXISTS topic CASCADE;
DROP TABLE IF EXISTS topics CASCADE;

-- Drop restaurant related tables (both singular and plural)
DROP TABLE IF EXISTS restaurant_users CASCADE;
DROP TABLE IF EXISTS restaurant_hours CASCADE;
DROP TABLE IF EXISTS restaurant CASCADE;
DROP TABLE IF EXISTS restaurants CASCADE;
