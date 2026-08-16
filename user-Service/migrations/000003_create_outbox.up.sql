CREATE TABLE IF NOT EXISTS outbox (
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    event_type TEXT    NOT NULL,
    payload    JSONB   NOT NULL CHECK (jsonb_typeof(payload) = 'object'),
    sent       BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_outbox_unsent ON outbox (created_at) WHERE sent = false;
