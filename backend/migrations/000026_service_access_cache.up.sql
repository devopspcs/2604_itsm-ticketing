CREATE TABLE IF NOT EXISTS service_access_cache (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id UUID NOT NULL REFERENCES external_services(id) ON DELETE CASCADE,
    service_name VARCHAR(255) NOT NULL,
    service_type VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    account_name VARCHAR(255) DEFAULT '',
    roles TEXT[] DEFAULT '{}',
    status VARCHAR(50) DEFAULT 'active',
    synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_service_access_cache_email ON service_access_cache(email);
CREATE INDEX idx_service_access_cache_service ON service_access_cache(service_id);

-- Track last sync time
CREATE TABLE IF NOT EXISTS sync_status (
    id VARCHAR(50) PRIMARY KEY,
    last_synced_at TIMESTAMPTZ,
    status VARCHAR(50) DEFAULT 'idle',
    total_records INT DEFAULT 0
);

INSERT INTO sync_status (id, status) VALUES ('service_access', 'never') ON CONFLICT DO NOTHING;
