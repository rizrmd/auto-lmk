CREATE TABLE showroom_settings (
    tenant_id INTEGER PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
    address TEXT,
    phone VARCHAR(20),
    email VARCHAR(255),
    business_hours TEXT,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    map_embed TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
