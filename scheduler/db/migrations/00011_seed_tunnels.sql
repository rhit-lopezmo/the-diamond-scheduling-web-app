-- +goose Up
INSERT INTO tunnels (name)
SELECT 'Tunnel ' || i::text
FROM generate_series(1,10) AS g(i);

-- +goose Down
-- Forward-only policy: no down migration provided.
