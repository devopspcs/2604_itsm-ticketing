-- Revert org structure hierarchy back to: Department (top) -> Division (child) -> Team (child)

-- Drop new constraints
ALTER TABLE teams DROP CONSTRAINT IF EXISTS teams_department_id_fkey;
ALTER TABLE departments DROP CONSTRAINT IF EXISTS departments_division_id_fkey;
ALTER TABLE departments DROP CONSTRAINT IF EXISTS departments_division_id_code_key;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_department_id_fkey;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_division_id_fkey;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_team_id_fkey;

-- Drop new indexes
DROP INDEX IF EXISTS idx_departments_division_id;
DROP INDEX IF EXISTS idx_teams_department_id;

-- Revert teams: add division_id, drop department_id
ALTER TABLE teams ADD COLUMN IF NOT EXISTS division_id UUID;
ALTER TABLE teams DROP COLUMN IF EXISTS department_id;

-- Revert departments: drop division_id
ALTER TABLE departments DROP COLUMN IF EXISTS division_id;
ALTER TABLE departments ADD CONSTRAINT departments_name_key UNIQUE (name);

-- Revert divisions: add department_id
ALTER TABLE divisions ADD COLUMN IF NOT EXISTS department_id UUID;

-- Re-add original constraints
ALTER TABLE divisions
    ADD CONSTRAINT divisions_department_id_fkey
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE RESTRICT;

ALTER TABLE teams
    ADD CONSTRAINT teams_division_id_fkey
    FOREIGN KEY (division_id) REFERENCES divisions(id) ON DELETE RESTRICT;

ALTER TABLE users
    ADD CONSTRAINT users_department_id_fkey
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE RESTRICT;

ALTER TABLE users
    ADD CONSTRAINT users_division_id_fkey
    FOREIGN KEY (division_id) REFERENCES divisions(id) ON DELETE RESTRICT;

ALTER TABLE users
    ADD CONSTRAINT users_team_id_fkey
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE RESTRICT;
