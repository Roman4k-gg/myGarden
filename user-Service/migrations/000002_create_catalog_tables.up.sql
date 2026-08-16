CREATE TABLE IF NOT EXISTS plants (
    id                     BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name                   TEXT NOT NULL,
    description            TEXT NOT NULL DEFAULT '',
    watering_interval_days INT  NOT NULL CHECK (watering_interval_days > 0),
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_plants_created_at ON plants (created_at);

CREATE TABLE IF NOT EXISTS user_favorites (
    user_id    BIGINT NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
    plant_id   BIGINT NOT NULL REFERENCES plants(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, plant_id)
);

CREATE INDEX IF NOT EXISTS idx_user_favorites_plant_id ON user_favorites (plant_id);

INSERT INTO plants (name, description, watering_interval_days) VALUES
    ('Фикус',    'Неприхотливое комнатное растение', 7),
    ('Монстера', 'Любит влажность и непрямой свет',  5),
    ('Кактус',   'Редкий полив, много солнца',       21)
ON CONFLICT DO NOTHING;
