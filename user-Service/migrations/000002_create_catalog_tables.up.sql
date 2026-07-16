CREATE TABLE IF NOT EXISTS plants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    watering_interval_days INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_favorites (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plant_id INT NOT NULL REFERENCES plants(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, plant_id)
);

INSERT INTO plants (name, description, watering_interval_days) VALUES 
('Фикус', 'Неприхотливое комнатное растение', 7),
('Монстера', 'Любит влажность и непрямой свет', 5),
('Кактус', 'Редкий полив, много солнца', 21)
ON CONFLICT DO NOTHING;
