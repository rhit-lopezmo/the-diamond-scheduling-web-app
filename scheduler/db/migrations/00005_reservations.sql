-- +goose Up
CREATE TABLE reservations (
    id                     uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    reservation_kind       reservation_kind NOT NULL,

    tunnel_id              int REFERENCES tunnels(id) ON DELETE RESTRICT,
    coach_id               uuid REFERENCES coaches(id) ON DELETE RESTRICT,

    customer_first_name    text NOT NULL,
    customer_last_name     text NOT NULL,
    customer_phone         text NOT NULL,
    customer_email         text,

    start_time             timestamptz NOT NULL,
    duration_minutes       int NOT NULL CHECK (duration_minutes > 0),
    end_time               timestamptz NOT NULL,

    status                 reservation_status NOT NULL DEFAULT 'held',
    notes                  text,

    created_at             timestamptz NOT NULL DEFAULT now(),
    updated_at             timestamptz NOT NULL DEFAULT now(),

    CHECK (reservation_kind <> 'tunnel' OR (tunnel_id IS NOT NULL AND coach_id IS NULL)),
    CHECK (reservation_kind <> 'lesson' OR (tunnel_id IS NOT NULL AND coach_id IS NOT NULL))
);
-- +goose Down
-- Forward-only policy: no down migration provided.
