-- Migration: Add agent role and update ticket types
-- 1. Update role constraint to include 'agent'
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role IN ('user', 'agent', 'approver', 'admin'));

-- 2. Update ticket type constraint to replace 'helpdesk_request' with 'request'
-- First, update existing helpdesk_request tickets to 'request'
UPDATE tickets SET type = 'request' WHERE type = 'helpdesk_request';

-- Drop old constraint and add new one
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS tickets_type_check;
ALTER TABLE tickets ADD CONSTRAINT tickets_type_check CHECK (type IN ('incident', 'request', 'change_request'));

-- 3. Migrate existing users with role 'user' who should be agents
-- (This is a manual step - admin should update specific users to 'agent' role via the UI)
-- For now, we just ensure the constraint allows the new role.
