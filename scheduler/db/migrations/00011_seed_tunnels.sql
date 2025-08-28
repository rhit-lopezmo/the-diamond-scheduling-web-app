-- +goose Up
INSERT INTO tunnels (name)
VALUES ('Tunnel 1'),
       ('Tunnel 2'),
       ('Tunnel 3'),
       ('Tunnel 4'),
       ('Tunnel 5'),
       ('Tunnel 6'),
       ('Tunnel 7'),
       ('Tunnel 8'),
       ('Tunnel 9');

-- +goose Down
-- Forward-only policy: no down migration provided.
