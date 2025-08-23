-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'reservation_kind') THEN
    CREATE TYPE reservation_kind AS ENUM ('tunnel', 'lesson');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'reservation_status') THEN
    CREATE TYPE reservation_status AS ENUM ('held', 'confirmed', 'cancelled', 'completed', 'no_show');
  END IF;
END$$;
-- +goose StatementEnd

-- +goose Down
-- Forward-only policy: no down migration provided.
