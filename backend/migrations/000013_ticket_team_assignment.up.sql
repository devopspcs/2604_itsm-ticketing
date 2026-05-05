-- Add email field to teams table for team notifications
ALTER TABLE teams ADD COLUMN email VARCHAR(255);

-- Add assigned_team_id to tickets for team-level assignment
ALTER TABLE tickets ADD COLUMN assigned_team_id UUID REFERENCES teams(id) ON DELETE SET NULL;

CREATE INDEX idx_tickets_assigned_team_id ON tickets(assigned_team_id);
