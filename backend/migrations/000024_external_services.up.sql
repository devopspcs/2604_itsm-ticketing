CREATE TABLE IF NOT EXISTS external_services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'pritunl', 'pose', 'custom_api'
    url TEXT NOT NULL,
    auth_username VARCHAR(255) DEFAULT '',
    auth_token TEXT DEFAULT '',
    auth_password TEXT DEFAULT '',
    extra_config JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed existing services
INSERT INTO external_services (name, type, url, auth_token, extra_config, is_active) VALUES
('POSe', 'pose', 'https://pose-api.pcsindonesia.co.id/master/user?filter=pcsindonesia.co.id&order=created_at:-1&page_size=-1&page=1', '', '{}', true)
ON CONFLICT DO NOTHING;
