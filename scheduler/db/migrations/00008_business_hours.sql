-- +goose Up
CREATE TABLE business_hours (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  dow        INT  NOT NULL CHECK (dow BETWEEN 0 AND 6), -- 0=Sun … 6=Sat
  open_time  TIME NOT NULL,
  close_time TIME NOT NULL,
  is_open    BOOLEAN NOT NULL DEFAULT TRUE,
  UNIQUE (dow),
  CHECK (open_time < close_time)
);

-- +goose Down
-- Forward-only policy: no down migration provided.
