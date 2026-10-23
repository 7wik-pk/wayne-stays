-- +goose Up
-- Baseline: mirrors the table GORM's AutoMigrate used to create, so databases
-- created before versioned migrations existed are adopted without changes.
CREATE TABLE IF NOT EXISTS accommodations (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    slug        text NOT NULL,
    title       text NOT NULL,
    description text,
    location    text,
    price_start bigint NOT NULL, -- in cents
    capacity    bigint DEFAULT 2,
    category    text,            -- urban-luxury, island, historic, countryside
    status      text DEFAULT 'active', -- active, maintenance, inactive
    amenities   jsonb,
    images      jsonb,
    created_at  timestamptz,
    updated_at  timestamptz,
    deleted_at  timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_accommodations_slug ON accommodations (slug);
CREATE INDEX IF NOT EXISTS idx_accommodations_deleted_at ON accommodations (deleted_at);
CREATE INDEX IF NOT EXISTS idx_accommodations_category ON accommodations (category);
CREATE INDEX IF NOT EXISTS idx_accommodations_status ON accommodations (status);

-- +goose Down
DROP TABLE IF EXISTS accommodations;
