CREATE TABLE IF NOT EXISTS car_photos (
    id SERIAL PRIMARY KEY,
    car_id INTEGER NOT NULL REFERENCES cars(id) ON DELETE CASCADE,
    file_path VARCHAR(500) NOT NULL,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_car_photos_car_id ON car_photos(car_id);
CREATE INDEX idx_car_photos_display_order ON car_photos(display_order);
