-- +goose Up
CREATE TABLE reservations (
  id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  reservation_kind     reservation_kind NOT NULL,

  tunnel_id            UUID REFERENCES tunnels(id) ON DELETE RESTRICT,
  coach_id             UUID REFERENCES coaches(id) ON DELETE RESTRICT,

  -- Inline customer info (MVP: no persistent customers yet)
  customer_first_name  TEXT NOT NULL,
  customer_last_name   TEXT NOT NULL,
  customer_phone       TEXT,
  customer_email       TEXT,

  -- Timing (UTC); end_time is generated and used for overlap checks
  start_time           TIMESTAMPTZ NOT NULL,
  duration_minutes     INT NOT NULL CHECK (duration_minutes > 0),
  end_time             TIMESTAMPTZ GENERATED ALWAYS AS
                         (start_time + make_interval(mins => duration_minutes)) STORED,

  -- Status & bookkeeping
  status               reservation_status NOT NULL DEFAULT 'held',
  source               TEXT NOT NULL DEFAULT 'web',  -- 'web'|'admin'
  is_custom_duration   BOOLEAN NOT NULL DEFAULT FALSE,
  price_cents          INT CHECK (price_cents IS NULL OR price_cents >= 0),
  notes                TEXT,

  -- Stripe placeholders
  stripe_payment_intent_id   TEXT,
  stripe_checkout_session_id TEXT,

  created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),

  -- Kind-specific presence rules
  CHECK (reservation_kind <> 'tunnel' OR (tunnel_id IS NOT NULL AND coach_id IS NULL)),
  CHECK (reservation_kind <> 'lesson' OR (tunnel_id IS NOT NULL AND coach_id IS NOT NULL))
);

-- +goose Down
-- Forward-only policy: no down migration provided.
