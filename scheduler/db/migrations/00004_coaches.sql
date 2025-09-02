-- +goose Up
CREATE TABLE coaches (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  first_name  TEXT NOT NULL,
  last_name   TEXT NOT NULL,
  email       TEXT,
  phone       TEXT NOT NULL,
  is_active   BOOLEAN NOT NULL DEFAULT TRUE,
  specialties coach_specialty[] NOT NULL DEFAULT '{}',
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
-- Forward-only policy: no down migration provided.
