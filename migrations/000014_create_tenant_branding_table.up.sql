CREATE TABLE tenant_branding (
    tenant_id INTEGER PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
    logo_path VARCHAR(255),
    favicon_path VARCHAR(255),
    custom_title VARCHAR(255),
    custom_subtitle VARCHAR(255),
    promo_text TEXT,
    header_style VARCHAR(50) DEFAULT 'default',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
