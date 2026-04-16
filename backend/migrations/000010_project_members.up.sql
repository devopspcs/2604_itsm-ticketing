-- Project members for invitation-based access
CREATE TABLE project_members (
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id),
    role       VARCHAR(20) NOT NULL DEFAULT 'member' CHECK (role IN ('owner', 'member')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (project_id, user_id)
);
CREATE INDEX idx_project_members_user ON project_members(user_id);

-- Seed: add existing project creators as owners
INSERT INTO project_members (project_id, user_id, role)
SELECT id, created_by, 'owner' FROM projects
ON CONFLICT DO NOTHING;
