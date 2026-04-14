DROP INDEX IF EXISTS idx_users_position;
DROP INDEX IF EXISTS idx_users_team_id;
DROP INDEX IF EXISTS idx_users_division_id;
DROP INDEX IF EXISTS idx_users_department_id;

ALTER TABLE users
    DROP COLUMN IF EXISTS position,
    DROP COLUMN IF EXISTS team_id,
    DROP COLUMN IF EXISTS division_id,
    DROP COLUMN IF EXISTS department_id;

DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS divisions;
DROP TABLE IF EXISTS departments;
