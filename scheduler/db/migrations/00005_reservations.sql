-- +goose Up
CREATE TABLE reservations (
    id                     uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    reservation_kind       reservation_kind NOT NULL,

    tunnel_id              uuid REFERENCES tunnels(id) ON DELETE RESTRICT,
    coach_id               uuid REFERENCES coaches(id) ON DELETE RESTRICT,

    customer_first_name    text NOT NULL,
    customer_last_name     text NOT NULL,
    customer_phone         text,
    customer_email         text,

    -- Keep timestamptz; we’ll compute end_time in app or trigger
    start_time             timestamptz NOT NULL,
    duration_minutes       int NOT NULL CHECK (duration_minutes > 0),
    end_time               timestamptz NOT NULL,

    status                 reservation_status NOT NULL DEFAULT 'held',
    source                 text NOT NULL DEFAULT 'web',  -- 'web'|'admin'
    is_custom_duration     boolean NOT NULL DEFAULT false,
    price_cents            int CHECK (price_cents IS NULL OR price_cents >= 0),
    notes                  text,

    stripe_payment_intent_id   text,
    stripe_checkout_session_id text,

    created_at             timestamptz NOT NULL DEFAULT now(),

    CHECK (reservation_kind <> 'tunnel' OR (tunnel_id IS NOT NULL AND coach_id IS NULL)),
    CHECK (reservation_kind <> 'lesson' OR (tunnel_id IS NOT NULL AND coach_id IS NOT NULL))
);
-- +goose Down
-- Forward-only policy: no down migration provided.
