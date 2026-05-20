CREATE TABLE applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    code VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    icon VARCHAR(50) DEFAULT 'apps',
    color VARCHAR(7) DEFAULT '#1976d2',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE user_app_access (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    app_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    granted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    granted_by UUID REFERENCES users(id),
    UNIQUE(user_id, app_id)
);

CREATE INDEX idx_user_app_access_user ON user_app_access(user_id);
CREATE INDEX idx_user_app_access_app ON user_app_access(app_id);

-- Seed default applications
INSERT INTO applications (name, code, description, icon, color) VALUES
('Ticketing System', 'ticketing', 'IT Service Management & Ticketing', 'confirmation_number', '#d32f2f'),
('Project Board', 'project-board', 'Project Management & Kanban Board', 'dashboard', '#1976d2');

-- Give all existing users access to both apps
INSERT INTO user_app_access (user_id, app_id, role)
SELECT u.id, a.id, u.role
FROM users u CROSS JOIN applications a;
