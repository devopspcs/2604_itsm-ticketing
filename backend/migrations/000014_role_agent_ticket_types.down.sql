-- Rollback: Remove agent role and revert ticket types
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role IN ('user', 'approver', 'admin'));

UPDATE tickets SET type = 'helpdesk_request' WHERE type = 'request';

ALTER TABLE tickets DROP CONSTRAINT IF EXISTS tickets_type_check;
ALTER TABLE tickets ADD CONSTRAINT tickets_type_check CHECK (type IN ('change_request', 'incident', 'helpdesk_request'));
