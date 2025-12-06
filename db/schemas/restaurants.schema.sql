-- =========================
-- RESTAURANTS
-- =========================
CREATE TABLE IF NOT EXISTS restaurants (
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
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by     UUID REFERENCES users(id),
    updated_by     UUID REFERENCES users(id),
    delete_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status         TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended'))
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


CREATE TABLE IF NOT EXISTS restaurant_users (
    id             INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    restaurant_id INT NOT NULL REFERENCES restaurant(id) ON DELETE CASCADE,
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role          TEXT NOT NULL CHECK (role IN ('owner', 'manager', 'staff')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    delete_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (restaurant_id, user_id)
);

CREATE TRIGGER trg_restaurant_users_updated_at
BEFORE UPDATE ON restaurant_users
FOR EACH ROW EXECUTE FUNCTION set_updated_at();
