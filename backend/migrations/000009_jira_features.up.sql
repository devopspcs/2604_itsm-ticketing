-- Issue Types (predefined)
CREATE TABLE issue_types (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL UNIQUE,
    icon        VARCHAR(50),
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Issue Type Schemes
CREATE TABLE issue_type_schemes (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_issue_type_schemes_project_id ON issue_type_schemes(project_id);

-- Issue Type Scheme Items
CREATE TABLE issue_type_scheme_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scheme_id       UUID NOT NULL REFERENCES issue_type_schemes(id) ON DELETE CASCADE,
    issue_type_id   UUID NOT NULL REFERENCES issue_types(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_issue_type_scheme_items_scheme_id ON issue_type_scheme_items(scheme_id);

-- Custom Fields
CREATE TABLE custom_fields (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    field_type  VARCHAR(50) NOT NULL,
    is_required BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_custom_fields_project_id ON custom_fields(project_id);

-- Custom Field Options
CREATE TABLE custom_field_options (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    field_id     UUID NOT NULL REFERENCES custom_fields(id) ON DELETE CASCADE,
    option_value VARCHAR(255) NOT NULL,
    option_order INTEGER NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_custom_field_options_field_id ON custom_field_options(field_id);

-- Custom Field Values
CREATE TABLE custom_field_values (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id  UUID NOT NULL REFERENCES project_records(id) ON DELETE CASCADE,
    field_id   UUID NOT NULL REFERENCES custom_fields(id) ON DELETE CASCADE,
    value      TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_custom_field_values_record_id ON custom_field_values(record_id);
CREATE INDEX idx_custom_field_values_field_id ON custom_field_values(field_id);

-- Workflows
CREATE TABLE workflows (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    initial_status  VARCHAR(100) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workflows_project_id ON workflows(project_id);

-- Workflow Statuses
CREATE TABLE workflow_statuses (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    status_name VARCHAR(100) NOT NULL,
    status_order INTEGER NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workflow_statuses_workflow_id ON workflow_statuses(workflow_id);

-- Workflow Transitions
CREATE TABLE workflow_transitions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id     UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    from_status_id  UUID NOT NULL REFERENCES workflow_statuses(id),
    to_status_id    UUID NOT NULL REFERENCES workflow_statuses(id),
    validation_rule TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workflow_transitions_workflow_id ON workflow_transitions(workflow_id);

-- Sprints
CREATE TABLE sprints (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id          UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name                VARCHAR(255) NOT NULL,
    goal                TEXT,
    start_date          DATE,
    end_date            DATE NOT NULL,
    status              VARCHAR(50) NOT NULL DEFAULT 'Planned',
    actual_start_date   DATE,
    actual_end_date     DATE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sprints_project_id ON sprints(project_id);
CREATE INDEX idx_sprints_status ON sprints(status);

-- Sprint Records
CREATE TABLE sprint_records (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sprint_id  UUID NOT NULL REFERENCES sprints(id) ON DELETE CASCADE,
    record_id  UUID NOT NULL REFERENCES project_records(id) ON DELETE CASCADE,
    priority   INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sprint_records_sprint_id ON sprint_records(sprint_id);
CREATE INDEX idx_sprint_records_record_id ON sprint_records(record_id);

-- Comments
CREATE TABLE comments (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id  UUID NOT NULL REFERENCES project_records(id) ON DELETE CASCADE,
    author_id  UUID NOT NULL REFERENCES users(id),
    text       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comments_record_id ON comments(record_id);
CREATE INDEX idx_comments_author_id ON comments(author_id);

-- Comment Mentions
CREATE TABLE comment_mentions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    comment_id          UUID NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    mentioned_user_id   UUID NOT NULL REFERENCES users(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comment_mentions_comment_id ON comment_mentions(comment_id);
CREATE INDEX idx_comment_mentions_mentioned_user_id ON comment_mentions(mentioned_user_id);

-- Attachments
CREATE TABLE attachments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id   UUID NOT NULL REFERENCES project_records(id) ON DELETE CASCADE,
    file_name   VARCHAR(500) NOT NULL,
    file_size   BIGINT NOT NULL,
    file_type   VARCHAR(100),
    file_path   VARCHAR(1000) NOT NULL,
    uploader_id UUID NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_attachments_record_id ON attachments(record_id);
CREATE INDEX idx_attachments_uploader_id ON attachments(uploader_id);

-- Labels
CREATE TABLE labels (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name       VARCHAR(100) NOT NULL,
    color      VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_labels_project_id ON labels(project_id);

-- Record Labels
CREATE TABLE record_labels (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id  UUID NOT NULL REFERENCES project_records(id) ON DELETE CASCADE,
    label_id   UUID NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_record_labels_record_id ON record_labels(record_id);
CREATE INDEX idx_record_labels_label_id ON record_labels(label_id);

-- Extend project_records table
ALTER TABLE project_records ADD COLUMN issue_type_id UUID REFERENCES issue_types(id);
ALTER TABLE project_records ADD COLUMN status VARCHAR(100);
ALTER TABLE project_records ADD COLUMN parent_record_id UUID REFERENCES project_records(id);

CREATE INDEX idx_project_records_issue_type_id ON project_records(issue_type_id);
CREATE INDEX idx_project_records_status ON project_records(status);
CREATE INDEX idx_project_records_parent_record_id ON project_records(parent_record_id);
