-- Rollback: remove all seeded org data
UPDATE users SET department_id = NULL, division_id = NULL, team_id = NULL;
DELETE FROM teams;
DELETE FROM departments;
DELETE FROM divisions;
