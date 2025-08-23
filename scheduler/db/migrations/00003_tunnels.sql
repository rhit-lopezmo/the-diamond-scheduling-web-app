-- +goose Up
CREATE TABLE tunnels (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name       TEXT NOT NULL,                -- e.g., "Tunnel 1"
  is_active  BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (name)
);

-- +goose Down
-- Forward-only policy: no down migration provided.
