-- +goose Up
CREATE INDEX idx_reservations_tunnel_start ON reservations (tunnel_id, start_time) WHERE tunnel_id IS NOT NULL;
CREATE INDEX idx_reservations_coach_start  ON reservations (coach_id, start_time)  WHERE coach_id  IS NOT NULL;

-- +goose Down
-- Forward-only policy: no down migration provided.
