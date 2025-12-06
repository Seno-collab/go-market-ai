ALTER TABLE "user" RENAME TO users;
ALTER TABLE "restaurant" RENAME TO restaurants;

CREATE TYPE restaurant_status AS ENUM ('active', 'inactive', 'suspended');

ALTER TABLE restaurants
DROP COLUMN IF EXISTS user_id;

ALTER TABLE restaurants
ADD COLUMN created_by UUID REFERENCES users(id),
ADD COLUMN status restaurant_status NOT NULL DEFAULT 'active',
ADD COLUMN updated_by UUID REFERENCES users(id);


CREATE TABLE IF NOT EXISTS restaurant_users (
    id            INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    restaurant_id INT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role          TEXT NOT NULL CHECK (role IN ('owner', 'manager', 'staff')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (restaurant_id, user_id)
);
