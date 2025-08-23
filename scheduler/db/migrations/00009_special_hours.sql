-- +goose Up
CREATE TABLE special_hours (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  on_date    DATE NOT NULL,                 -- local calendar date
  open_time  TIME NOT NULL,
  close_time TIME NOT NULL,
  is_open    BOOLEAN NOT NULL DEFAULT TRUE, -- false => closed all day
  notes      TEXT,
  UNIQUE (on_date),
  CHECK (open_time < close_time)
);

-- +goose Down
-- Forward-only policy: no down migration provided.
