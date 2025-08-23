-- +goose Up
-- Local wall times for America/Indiana/Indianapolis.
INSERT INTO business_hours (dow, open_time, close_time, is_open) VALUES
  (0, '10:00', '17:00', TRUE),  -- Sun
  (1, '15:00', '21:00', TRUE),  -- Mon
  (2, '15:00', '21:00', TRUE),  -- Tue
  (3, '15:00', '21:00', TRUE),  -- Wed
  (4, '15:00', '21:00', TRUE),  -- Thu
  (5, '15:00', '21:00', TRUE),  -- Fri
  (6, '10:00', '17:00', TRUE)   -- Sat
ON CONFLICT (dow) DO UPDATE
SET open_time = EXCLUDED.open_time,
    close_time = EXCLUDED.close_time,
    is_open    = EXCLUDED.is_open;

-- +goose Down
-- Forward-only policy: no down migration provided.
