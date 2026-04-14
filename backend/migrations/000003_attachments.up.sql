CREATE TABLE ticket_attachments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id   UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    uploaded_by UUID NOT NULL REFERENCES users(id),
    filename    VARCHAR(255) NOT NULL,
    file_size   BIGINT NOT NULL DEFAULT 0,
    mime_type   VARCHAR(100) NOT NULL DEFAULT 'application/octet-stream',
    file_data   BYTEA NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_attachments_ticket_id ON ticket_attachments(ticket_id);
