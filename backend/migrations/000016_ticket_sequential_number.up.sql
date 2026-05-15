-- Table to track sequential counters per ticket type
CREATE TABLE ticket_sequences (
    ticket_type VARCHAR(50) PRIMARY KEY,
    last_number BIGINT NOT NULL DEFAULT 0
);

-- Seed initial counters for each ticket type
INSERT INTO ticket_sequences (ticket_type, last_number) VALUES
    ('incident', 0),
    ('request', 0),
    ('change_request', 0);

-- Add ticket_number column to tickets table
ALTER TABLE tickets ADD COLUMN ticket_number VARCHAR(20);

-- Create unique index on ticket_number
CREATE UNIQUE INDEX idx_tickets_ticket_number ON tickets(ticket_number);

-- Backfill existing tickets with sequential numbers based on created_at order
DO $$
DECLARE
    rec RECORD;
    counter_inc BIGINT := 0;
    counter_req BIGINT := 0;
    counter_chg BIGINT := 0;
    new_number VARCHAR(20);
BEGIN
    FOR rec IN SELECT id, type FROM tickets ORDER BY created_at ASC
    LOOP
        CASE rec.type
            WHEN 'incident' THEN
                counter_inc := counter_inc + 1;
                new_number := 'INC-' || LPAD(counter_inc::TEXT, 6, '0');
            WHEN 'request' THEN
                counter_req := counter_req + 1;
                new_number := 'REQ-' || LPAD(counter_req::TEXT, 6, '0');
            WHEN 'change_request' THEN
                counter_chg := counter_chg + 1;
                new_number := 'CHG-' || LPAD(counter_chg::TEXT, 6, '0');
            ELSE
                counter_req := counter_req + 1;
                new_number := 'REQ-' || LPAD(counter_req::TEXT, 6, '0');
        END CASE;

        UPDATE tickets SET ticket_number = new_number WHERE id = rec.id;
    END LOOP;

    -- Update the sequence counters
    UPDATE ticket_sequences SET last_number = counter_inc WHERE ticket_type = 'incident';
    UPDATE ticket_sequences SET last_number = counter_req WHERE ticket_type = 'request';
    UPDATE ticket_sequences SET last_number = counter_chg WHERE ticket_type = 'change_request';
END $$;

-- Make ticket_number NOT NULL after backfill
ALTER TABLE tickets ALTER COLUMN ticket_number SET NOT NULL;
