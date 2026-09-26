-- Simple Commerce revamp: customer auth and a read-only product catalog.
-- PostgreSQL 14+. Fresh database schema; do not apply to an existing database
-- without a migration/data-preservation plan.

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    full_name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'disabled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT users_email_not_blank CHECK (length(btrim(email)) > 0),
    CONSTRAINT users_full_name_not_blank CHECK (length(btrim(full_name)) > 0)
);
-- Normalize email (trim + lower) in the application before insert and login;
-- this index also enforces case-insensitive uniqueness at the database boundary.
CREATE UNIQUE INDEX users_email_lower_uniq ON users (lower(email));

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT refresh_tokens_revoked_order CHECK (revoked_at IS NULL OR revoked_at >= created_at)
);
CREATE INDEX refresh_tokens_user_id_idx ON refresh_tokens (user_id);
-- Store a hash, never a raw refresh token. Logout revokes the presented session.

CREATE TABLE categories (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT categories_name_not_blank CHECK (length(btrim(name)) > 0),
    CONSTRAINT categories_slug_not_blank CHECK (length(btrim(slug)) > 0)
);

CREATE TABLE products (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    sku TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    price NUMERIC(15, 2) NOT NULL CHECK (price >= 0),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT products_sku_not_blank CHECK (length(btrim(sku)) > 0),
    CONSTRAINT products_name_not_blank CHECK (length(btrim(name)) > 0)
);
-- For active products filtered by category, sorted newest then id; revisit
-- after EXPLAIN ANALYZE with the actual list query and seeded data.
CREATE INDEX products_browse_category_idx
    ON products (category_id, created_at DESC, id DESC)
    WHERE is_active = true;
CREATE INDEX products_browse_idx
    ON products (created_at DESC, id DESC)
    WHERE is_active = true;

INSERT INTO categories (name, slug)
VALUES ('Electronics', 'electronics'), ('Books', 'books'), ('Home', 'home');

INSERT INTO products (category_id, sku, name, description, price, is_active, created_at, updated_at)
SELECT
    1 + ((n - 1) % 3),
    'SKU-' || lpad(n::text, 6, '0'),
    'Benchmark Product ' || n,
    'Deterministic catalog fixture ' || n,
    ((n * 7919) % 500000 + 1000)::numeric(15, 2),
    true,
    TIMESTAMPTZ '2026-01-01 00:00:00+00' + (n || ' seconds')::interval,
    TIMESTAMPTZ '2026-01-01 00:00:00+00' + (n || ' seconds')::interval
FROM generate_series(1, 10000) AS n;
