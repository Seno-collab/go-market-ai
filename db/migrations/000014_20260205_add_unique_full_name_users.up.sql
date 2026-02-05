-- Deduplicate: append _1, _2, ... to duplicate full_name rows (keep the oldest)
WITH duplicates AS (
    SELECT id, full_name,
           ROW_NUMBER() OVER (PARTITION BY full_name ORDER BY created_at ASC) AS rn
    FROM "users"
)
UPDATE "users"
SET full_name = "users".full_name || '_' || (duplicates.rn - 1)
FROM duplicates
WHERE "users".id = duplicates.id AND duplicates.rn > 1;

ALTER TABLE "users" ADD CONSTRAINT users_full_name_unique UNIQUE (full_name);
