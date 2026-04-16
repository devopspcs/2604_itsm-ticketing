-- Multi-assign for project records
CREATE TABLE project_record_assignees (
    record_id  UUID NOT NULL REFERENCES project_records(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (record_id, user_id)
);
CREATE INDEX idx_project_record_assignees_user ON project_record_assignees(user_id);
