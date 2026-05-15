-- Add drive_file_id column for Google Drive storage
ALTER TABLE ticket_attachments ADD COLUMN drive_file_id VARCHAR(255);

-- Make file_data nullable (not required when using Drive storage)
ALTER TABLE ticket_attachments ALTER COLUMN file_data DROP NOT NULL;
