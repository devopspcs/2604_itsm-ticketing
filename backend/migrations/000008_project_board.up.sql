-- Projects
CREATE TABLE projects (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    icon_color  VARCHAR(50) DEFAULT '#3b82f6',
    created_by  UUID NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_projects_created_by ON projects(created_by);

-- Project Columns
CREATE TABLE project_columns (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    position    INTEGER NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_project_columns_project_id ON project_columns(project_id);

-- Project Records
CREATE TABLE project_records (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    column_id    UUID NOT NULL REFERENCES project_columns(id),
    project_id   UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title        VARCHAR(500) NOT NULL,
    description  TEXT DEFAULT '',
    assigned_to  UUID REFERENCES users(id),
    due_date     DATE,
    position     INTEGER NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT false,
    completed_at TIMESTAMPTZ,
    created_by   UUID NOT NULL REFERENCES users(id),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_project_records_column_id   ON project_records(column_id);
CREATE INDEX idx_project_records_project_id  ON project_records(project_id);
CREATE INDEX idx_project_records_assigned_to ON project_records(assigned_to);
CREATE INDEX idx_project_records_due_date    ON project_records(due_date);

-- Project Activity Logs
CREATE TABLE project_activity_logs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    record_id   UUID REFERENCES project_records(id) ON DELETE SET NULL,
    actor_id    UUID NOT NULL REFERENCES users(id),
    action      VARCHAR(100) NOT NULL,
    detail      TEXT DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_project_activity_logs_project_id ON project_activity_logs(project_id);
CREATE INDEX idx_project_activity_logs_actor_id   ON project_activity_logs(actor_id);
CREATE INDEX idx_project_activity_logs_created_at ON project_activity_logs(created_at);
