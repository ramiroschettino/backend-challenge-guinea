CREATE TABLE IF NOT EXISTS users_read (
    id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    display_name VARCHAR(255),
    tenant_id VARCHAR(100) NOT NULL,
    created_at VARCHAR(50) NOT NULL,
    
    PRIMARY KEY (id, tenant_id)
);

CREATE INDEX idx_users_read_tenant ON users_read(tenant_id);
CREATE INDEX idx_users_read_email ON users_read(email);
CREATE INDEX idx_users_read_created_at ON users_read(created_at);