-- Remove drive_file_id column
ALTER TABLE ticket_attachments DROP COLUMN IF EXISTS drive_file_id;

-- Restore file_data NOT NULL constraint
ALTER TABLE ticket_attachments ALTER COLUMN file_data SET NOT NULL;
