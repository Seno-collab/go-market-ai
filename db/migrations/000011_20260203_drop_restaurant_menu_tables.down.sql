-- Rollback: recreate all tables from initial migration
-- This is a copy of the relevant parts from 000001_20251116_024550_create-table.up.sql

-- =========================
-- RESTAURANTS
-- =========================
CREATE TABLE IF NOT EXISTS restaurant (
  id             INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name           TEXT NOT NULL,
  description    TEXT,
  address        TEXT,
  category       TEXT,
  city           TEXT,
  district       TEXT,
  logo_url       TEXT,
  banner_url     TEXT,
  phone_number   TEXT,
  website_url    TEXT,
  email          CITEXT,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_restaurant_updated_at
BEFORE UPDATE ON restaurant
FOR EACH ROW EXECUTE FUNCTION set_updated_at();


CREATE TABLE restaurant_hours (
  restaurant_id INT REFERENCES restaurant(id) ON DELETE CASCADE,
  day_of_week   INT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6), -- 0=Sun
  open_time     TIME,
  close_time    TIME,
  is_closed     BOOLEAN NOT NULL DEFAULT FALSE,
  PRIMARY KEY (restaurant_id, day_of_week)
);


-- ====== TOPIC (category cha-con) ======
CREATE TABLE IF NOT EXISTS topic (
  id            BIGSERIAL PRIMARY KEY,
  restaurant_id INT NOT NULL REFERENCES restaurant(id) ON DELETE CASCADE,
  name          TEXT NOT NULL,
  slug          TEXT,
  parent_id     BIGINT REFERENCES topic(id) ON DELETE CASCADE,
  sort_order    INT NOT NULL DEFAULT 0,
  is_active     BOOLEAN NOT NULL DEFAULT TRUE,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (restaurant_id, slug)
);
CREATE INDEX IF NOT EXISTS idx_topic_parent ON topic(parent_id);
CREATE TRIGGER trg_topic_updated_at
BEFORE UPDATE ON topic
FOR EACH ROW EXECUTE FUNCTION set_updated_at();



-- ====== ENUMS ======
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'menu_item_type') THEN
    CREATE TYPE menu_item_type AS ENUM ('dish', 'extra', 'beverage', 'combo');
  END IF;
END $$;


-- ====== MENU ITEM (món ăn / đồ uống / extra / combo) ======
CREATE TABLE IF NOT EXISTS menu_item (
  id             BIGSERIAL PRIMARY KEY,
  restaurant_id  INT NOT NULL REFERENCES restaurant(id) ON DELETE CASCADE,
  topic_id       BIGINT REFERENCES topic(id) ON DELETE SET NULL,
  type           menu_item_type NOT NULL DEFAULT 'dish',
  name           TEXT NOT NULL,
  description    TEXT,
  image_url      TEXT,
  sku            TEXT,
  base_price     NUMERIC(12,2) NOT NULL DEFAULT 0,    -- giá cơ bản (VND)
  is_active      BOOLEAN NOT NULL DEFAULT TRUE,
  sort_order     INT NOT NULL DEFAULT 0,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (base_price >= 0)
);
CREATE INDEX IF NOT EXISTS idx_menu_item_restaurant ON menu_item(restaurant_id);
CREATE INDEX IF NOT EXISTS idx_menu_item_topic ON menu_item(topic_id);
CREATE INDEX IF NOT EXISTS idx_menu_item_type ON menu_item(type);
CREATE TRIGGER trg_menu_item_updated_at
BEFORE UPDATE ON menu_item
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ====== VARIANT (size/kích cỡ… nếu cần) ======
CREATE TABLE IF NOT EXISTS menu_item_variant (
  id             BIGSERIAL PRIMARY KEY,
  menu_item_id   BIGINT NOT NULL REFERENCES menu_item(id) ON DELETE CASCADE,
  name           TEXT NOT NULL,                 -- ví dụ: "Size lớn"
  price_delta    NUMERIC(12,2) NOT NULL DEFAULT 0,
  is_default     BOOLEAN NOT NULL DEFAULT FALSE,
  sort_order     INT NOT NULL DEFAULT 0,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (price_delta >= 0)
);
CREATE INDEX IF NOT EXISTS idx_variant_item ON menu_item_variant(menu_item_id);
CREATE TRIGGER trg_variant_updated_at
BEFORE UPDATE ON menu_item_variant
FOR EACH ROW EXECUTE FUNCTION set_updated_at();




-- ====== OPTION GROUPS (nhóm lựa chọn: ví dụ "Chọn thêm", "Thêm thịt") ======
-- Gắn nhóm vào 1 hoặc nhiều món: qua bảng nối menu_item_option_group
CREATE TABLE IF NOT EXISTS option_group (
  id             BIGSERIAL PRIMARY KEY,
  restaurant_id  INT NOT NULL REFERENCES restaurant(id) ON DELETE CASCADE,
  name           TEXT NOT NULL,                 -- ví dụ: "Đồ thêm", "Chọn món chính"
  min_select     INT NOT NULL DEFAULT 0,        -- ràng buộc min/max
  max_select     INT,                           -- NULL = không giới hạn
  is_required    BOOLEAN NOT NULL DEFAULT FALSE,
  sort_order     INT NOT NULL DEFAULT 0,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (min_select >= 0),
  CHECK (max_select IS NULL OR max_select >= 0)
);
CREATE TRIGGER trg_option_group_updated_at
BEFORE UPDATE ON option_group
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE IF NOT EXISTS menu_item_option_group (
  menu_item_id  BIGINT NOT NULL REFERENCES menu_item(id) ON DELETE CASCADE,
  option_group_id BIGINT NOT NULL REFERENCES option_group(id) ON DELETE CASCADE,
  sort_order    INT NOT NULL DEFAULT 0,
  PRIMARY KEY (menu_item_id, option_group_id)
);

-- ====== OPTIONS (mục lựa chọn trong nhóm)
CREATE TABLE IF NOT EXISTS option_item (
  id               BIGSERIAL PRIMARY KEY,
  option_group_id  BIGINT NOT NULL REFERENCES option_group(id) ON DELETE CASCADE,
  name             TEXT,                        -- dùng khi không link tới menu_item extra
  linked_menu_item BIGINT REFERENCES menu_item(id) ON DELETE SET NULL, -- dùng khi tái sử dụng item 'extra'
  price_delta      NUMERIC(12,2) NOT NULL DEFAULT 0, -- phụ thu thêm (cộng vào base)
  quantity_min     INT NOT NULL DEFAULT 0,
  quantity_max     INT,                         -- NULL = không giới hạn
  sort_order       INT NOT NULL DEFAULT 0,
  is_active        BOOLEAN NOT NULL DEFAULT TRUE,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (price_delta >= 0),
  CHECK (quantity_min >= 0),
  CHECK (quantity_max IS NULL OR quantity_max >= 0)
);
CREATE INDEX IF NOT EXISTS idx_option_group ON option_item(option_group_id);
CREATE TRIGGER trg_option_item_updated_at
BEFORE UPDATE ON option_item
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ====== COMBO ======
CREATE TABLE IF NOT EXISTS combo_group (
  id             BIGSERIAL PRIMARY KEY,
  combo_item_id  BIGINT NOT NULL REFERENCES menu_item(id) ON DELETE CASCADE,
  name           TEXT NOT NULL,                 -- ví dụ: "Món chính", "Món phụ", "Đồ uống"
  min_select     INT NOT NULL DEFAULT 1,        -- VD: "chọn 1"
  max_select     INT NOT NULL DEFAULT 1,        -- VD: "chọn 1", hoặc "chọn 2"
  sort_order     INT NOT NULL DEFAULT 0,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (min_select >= 0),
  CHECK (max_select >= 0),
  CHECK (max_select >= min_select)
);
CREATE INDEX IF NOT EXISTS idx_combo_group_combo ON combo_group(combo_item_id);
CREATE TRIGGER trg_combo_group_updated_at
BEFORE UPDATE ON combo_group
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE IF NOT EXISTS combo_group_item (
  id               BIGSERIAL PRIMARY KEY,
  combo_group_id   BIGINT NOT NULL REFERENCES combo_group(id) ON DELETE CASCADE,
  menu_item_id     BIGINT NOT NULL REFERENCES menu_item(id) ON DELETE RESTRICT,
  price_delta      NUMERIC(12,2) NOT NULL DEFAULT 0,
  quantity_default INT NOT NULL DEFAULT 1,
  quantity_min     INT NOT NULL DEFAULT 0,
  quantity_max     INT,
  sort_order       INT NOT NULL DEFAULT 0,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (price_delta >= 0),
  CHECK (quantity_default >= 0),
  CHECK (quantity_min >= 0),
  CHECK (quantity_max IS NULL OR quantity_max >= 0)
);
CREATE INDEX IF NOT EXISTS idx_combo_group_item_group ON combo_group_item(combo_group_id);
CREATE TRIGGER trg_combo_group_item_updated_at
BEFORE UPDATE ON combo_group_item
FOR EACH ROW EXECUTE FUNCTION set_updated_at();
