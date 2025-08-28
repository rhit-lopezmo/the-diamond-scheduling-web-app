-- +goose Up
CREATE TABLE tunnels (
  id         int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name       text NOT NULL,                -- e.g., "Tunnel 1"
  is_active  boolean NOT NULL DEFAULT TRUE,
  UNIQUE (name)
);

-- +goose Down
-- Forward-only policy: no down migration provided.
