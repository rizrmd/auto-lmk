CREATE TABLE IF NOT EXISTS cars (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INTEGER NOT NULL,
    price BIGINT NOT NULL,
    mileage INTEGER,
    transmission VARCHAR(50),
    fuel_type VARCHAR(50),
    engine_cc INTEGER,
    seats INTEGER,
    color VARCHAR(50),
    description TEXT,
    status VARCHAR(50) DEFAULT 'available',
    is_featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_cars_tenant_id ON cars(tenant_id);
CREATE INDEX idx_cars_status ON cars(status);
CREATE INDEX idx_cars_brand ON cars(brand);
CREATE INDEX idx_cars_price ON cars(price);
CREATE INDEX idx_cars_year ON cars(year);
