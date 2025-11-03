
CREATE TABLE IF NOT EXISTS idempotency_keys (
    key VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(100) NOT NULL,
    result TEXT NOT NULL,           
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (key, tenant_id)
);

CREATE INDEX idx_idempotency_tenant ON idempotency_keys(tenant_id);

CREATE INDEX idx_idempotency_created_at ON idempotency_keys(created_at);