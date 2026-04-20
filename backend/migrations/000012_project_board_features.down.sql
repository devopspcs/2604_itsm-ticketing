-- 000012_project_board_features.down.sql

ALTER TABLE project_records DROP COLUMN IF EXISTS component_id;
DROP TABLE IF EXISTS release_records;
DROP TABLE IF EXISTS releases;
DROP TABLE IF EXISTS components;
