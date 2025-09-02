-- +goose Up
CREATE TABLE blackout_windows (
  id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  starts_at TIMESTAMPTZ NOT NULL, -- UTC
  ends_at   TIMESTAMPTZ NOT NULL, -- UTC
  reason    TEXT NOT NULL,
  CHECK (starts_at < ends_at),
  EXCLUDE USING gist (tstzrange(starts_at, ends_at, '[)') WITH &&)
);

-- +goose Down
-- Forward-only policy: no down migration provided.
