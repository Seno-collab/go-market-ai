ALTER TABLE restaurants
ADD COLUMN deleted_at TIMESTAMPTZ;

ALTER TABLE restaurant_users
ADD COLUMN deleted_at TIMESTAMPTZ;

ALTER TABLE "role" RENAME TO roles;
