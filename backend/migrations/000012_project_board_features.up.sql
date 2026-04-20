-- 000012_project_board_features.up.sql

-- Releases table
CREATE TABLE IF NOT EXISTS releases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    version TEXT NOT NULL,
    description TEXT DEFAULT '',
    start_date TIMESTAMPTZ,
    release_date TIMESTAMPTZ,
    status TEXT NOT NULL DEFAULT 'Planning'
        CHECK (status IN ('Planning', 'In Progress', 'Released', 'Archived')),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_releases_project_id ON releases(project_id);

-- Release-Record association table
CREATE TABLE IF NOT EXISTS release_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    release_id UUID NOT NULL REFERENCES releases(id) ON DELETE CASCADE,
    record_id UUID NOT NULL REFERENCES project_records(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(release_id, record_id)
);

CREATE INDEX idx_release_records_release_id ON release_records(release_id);
CREATE INDEX idx_release_records_record_id ON release_records(record_id);

-- Components table
CREATE TABLE IF NOT EXISTS components (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    lead_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_components_project_id ON components(project_id);

-- Add component_id to project_records
ALTER TABLE project_records
    ADD COLUMN IF NOT EXISTS component_id UUID REFERENCES components(id) ON DELETE SET NULL;

CREATE INDEX idx_project_records_component_id ON project_records(component_id);
