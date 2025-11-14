CREATE TABLE IF NOT EXISTS car_specs (
    id SERIAL PRIMARY KEY,
    car_id INTEGER NOT NULL REFERENCES cars(id) ON DELETE CASCADE,
    key VARCHAR(100) NOT NULL,
    value TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_car_specs_car_id ON car_specs(car_id);
CREATE INDEX idx_car_specs_key ON car_specs(key);
