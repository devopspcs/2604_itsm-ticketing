-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- users
CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name   VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    role        VARCHAR(50)  NOT NULL CHECK (role IN ('user','approver','admin')),
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- refresh_tokens
CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  VARCHAR(255) NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked     BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- tickets
CREATE TABLE tickets (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       VARCHAR(500) NOT NULL,
    description TEXT NOT NULL,
    type        VARCHAR(50)  NOT NULL CHECK (type IN ('change_request','incident','helpdesk_request')),
    category    VARCHAR(100),
    priority    VARCHAR(20)  NOT NULL CHECK (priority IN ('low','medium','high','critical')),
    status      VARCHAR(30)  NOT NULL DEFAULT 'open'
                    CHECK (status IN ('open','in_progress','pending_approval','approved','rejected','done')),
    created_by  UUID NOT NULL REFERENCES users(id),
    assigned_to UUID REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tickets_status      ON tickets(status);
CREATE INDEX idx_tickets_type        ON tickets(type);
CREATE INDEX idx_tickets_created_by  ON tickets(created_by);
CREATE INDEX idx_tickets_assigned_to ON tickets(assigned_to);
CREATE INDEX idx_tickets_created_at  ON tickets(created_at);

-- approval_configs
CREATE TABLE approval_configs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_type VARCHAR(50) NOT NULL,
    level       INT NOT NULL,
    approver_id UUID NOT NULL REFERENCES users(id),
    UNIQUE (ticket_type, level, approver_id)
);

-- approvals
CREATE TABLE approvals (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id   UUID NOT NULL REFERENCES tickets(id),
    approver_id UUID NOT NULL REFERENCES users(id),
    level       INT NOT NULL,
    decision    VARCHAR(20) CHECK (decision IN ('approved','rejected')),
    comment     TEXT,
    decided_at  TIMESTAMPTZ
);

-- activity_logs
CREATE TABLE activity_logs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id   UUID NOT NULL REFERENCES tickets(id),
    actor_id    UUID NOT NULL REFERENCES users(id),
    action      VARCHAR(50) NOT NULL,
    old_value   TEXT,
    new_value   TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_activity_logs_ticket_id ON activity_logs(ticket_id);

-- notifications
CREATE TABLE notifications (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id),
    ticket_id   UUID NOT NULL REFERENCES tickets(id),
    message     TEXT NOT NULL,
    is_read     BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);

-- webhook_configs
CREATE TABLE webhook_configs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url         TEXT NOT NULL,
    events      TEXT[] NOT NULL,
    secret_key  VARCHAR(255) NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- webhook_logs
CREATE TABLE webhook_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    webhook_id      UUID NOT NULL REFERENCES webhook_configs(id),
    event           VARCHAR(50) NOT NULL,
    payload         JSONB NOT NULL,
    response_status INT,
    attempt         INT NOT NULL DEFAULT 1,
    sent_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
