CREATE TABLE IF NOT EXISTS users_write (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    display_name VARCHAR(255),
    tenant_id VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT unique_email_per_tenant UNIQUE (email, tenant_id)
);

CREATE INDEX idx_users_write_tenant ON users_write(tenant_id);
CREATE INDEX idx_users_write_email ON users_write(email);