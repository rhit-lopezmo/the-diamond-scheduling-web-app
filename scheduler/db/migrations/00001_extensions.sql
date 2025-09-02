-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;   -- gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS btree_gist; -- GiST ops for = with ranges

-- +goose Down
-- Forward-only policy: no down migration provided.
