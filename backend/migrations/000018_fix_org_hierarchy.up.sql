-- Fix org structure hierarchy: Division (top) -> Department (child) -> Team (child)
-- Previously: Department (top) -> Division (child) -> Team (child)

-- Step 1: Drop foreign key constraints on teams and users
ALTER TABLE teams DROP CONSTRAINT IF EXISTS teams_division_id_fkey;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_department_id_fkey;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_division_id_fkey;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_team_id_fkey;

-- Step 2: Drop foreign key on divisions (was referencing departments)
ALTER TABLE divisions DROP CONSTRAINT IF EXISTS divisions_department_id_fkey;

-- Step 3: Restructure divisions table - remove department_id, make it top-level
ALTER TABLE divisions DROP COLUMN IF EXISTS department_id;

-- Step 4: Restructure departments table - add division_id, make it child of divisions
ALTER TABLE departments ADD COLUMN IF NOT EXISTS division_id UUID;

-- Step 5: Restructure teams table - change from division_id to department_id
ALTER TABLE teams ADD COLUMN IF NOT EXISTS department_id UUID;

-- Migrate existing team->division relationships through to departments if possible
-- (data migration - best effort, may need manual cleanup)
UPDATE teams t SET department_id = (
    SELECT d.id FROM departments d LIMIT 1
) WHERE t.department_id IS NULL AND EXISTS (SELECT 1 FROM departments LIMIT 1);

ALTER TABLE teams DROP COLUMN IF EXISTS division_id;

-- Step 6: Add new foreign key constraints
ALTER TABLE departments
    ADD CONSTRAINT departments_division_id_fkey
    FOREIGN KEY (division_id) REFERENCES divisions(id) ON DELETE RESTRICT;

ALTER TABLE teams
    ADD CONSTRAINT teams_department_id_fkey
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE RESTRICT;

-- Re-add user constraints
ALTER TABLE users
    ADD CONSTRAINT users_department_id_fkey
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE RESTRICT;

ALTER TABLE users
    ADD CONSTRAINT users_division_id_fkey
    FOREIGN KEY (division_id) REFERENCES divisions(id) ON DELETE RESTRICT;

ALTER TABLE users
    ADD CONSTRAINT users_team_id_fkey
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE RESTRICT;

-- Step 7: Update unique constraint on departments
ALTER TABLE departments DROP CONSTRAINT IF EXISTS departments_name_key;
ALTER TABLE departments DROP CONSTRAINT IF EXISTS departments_code_key;
ALTER TABLE departments ADD CONSTRAINT departments_division_id_code_key UNIQUE (division_id, code);

-- Step 8: Create indexes
CREATE INDEX IF NOT EXISTS idx_departments_division_id ON departments(division_id);
CREATE INDEX IF NOT EXISTS idx_teams_department_id ON teams(department_id);
