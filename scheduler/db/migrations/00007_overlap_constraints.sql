-- +goose Up
ALTER TABLE reservations
  ADD CONSTRAINT tunnel_no_overlap
  EXCLUDE USING gist (
    tunnel_id WITH =,
    tstzrange(start_time, end_time, '[)') WITH &&
  )
  WHERE (tunnel_id IS NOT NULL AND status IN ('held','confirmed'));

ALTER TABLE reservations
  ADD CONSTRAINT coach_no_overlap
  EXCLUDE USING gist (
    coach_id WITH =,
    tstzrange(start_time, end_time, '[)') WITH &&
  )
  WHERE (coach_id IS NOT NULL AND status IN ('held','confirmed'));

-- +goose Down
-- Forward-only policy: no down migration provided.
