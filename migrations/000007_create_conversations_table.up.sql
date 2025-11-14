CREATE TABLE IF NOT EXISTS conversations (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    sender_phone VARCHAR(50) NOT NULL,
    is_sales BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_conversations_tenant_id ON conversations(tenant_id);
CREATE INDEX idx_conversations_sender_phone ON conversations(sender_phone);
CREATE INDEX idx_conversations_created_at ON conversations(created_at DESC);
