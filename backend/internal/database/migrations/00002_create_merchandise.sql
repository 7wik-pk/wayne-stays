-- +goose Up
-- GORM's AutoMigrate pluralised this table to "merchandises". Rename it (and its
-- indexes) to match schema.dbml when upgrading a database created that way.
-- +goose StatementBegin
DO $$
BEGIN
    IF to_regclass('public.merchandises') IS NOT NULL
       AND to_regclass('public.merchandise') IS NULL THEN
        ALTER TABLE merchandises RENAME TO merchandise;
        ALTER INDEX IF EXISTS merchandises_pkey RENAME TO merchandise_pkey;
        ALTER INDEX IF EXISTS idx_merchandises_category RENAME TO idx_merchandise_category;
        ALTER INDEX IF EXISTS idx_merchandises_deleted_at RENAME TO idx_merchandise_deleted_at;
    END IF;
END
$$;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS merchandise (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name        text NOT NULL,
    description text,
    price       bigint NOT NULL, -- in cents
    category    text,            -- apparel, accessories, collectibles, home, lifestyle
    image_url   text,
    is_new      boolean DEFAULT false,
    created_at  timestamptz,
    updated_at  timestamptz,
    deleted_at  timestamptz
);

-- Categories were originally seeded capitalised ("Apparel").
UPDATE merchandise SET category = lower(category) WHERE category <> lower(category);

CREATE INDEX IF NOT EXISTS idx_merchandise_category ON merchandise (category);
CREATE INDEX IF NOT EXISTS idx_merchandise_is_new ON merchandise (is_new);
CREATE INDEX IF NOT EXISTS idx_merchandise_deleted_at ON merchandise (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS merchandise;
