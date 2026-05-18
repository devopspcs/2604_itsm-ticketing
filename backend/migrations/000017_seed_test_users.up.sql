-- Seed test users for all teams
-- Password: adminpcs (bcrypt hash)
-- Run: docker exec itsm-postgres-prod psql -U itsm -d itsm -f /path/to/this/file
-- Or copy-paste the content into: docker exec itsm-postgres-prod psql -U itsm -d itsm

-- Fix position constraint first
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_position_check;
ALTER TABLE users ADD CONSTRAINT users_position_check CHECK (position IN ('manager', 'leader', 'staff'));
UPDATE users SET position = 'manager' WHERE position = 'division_manager';

DO $$
DECLARE
  pw TEXT := '$2a$12$SQbU9JeG5UlGmoBx3yUaZODfA00D2g8t28HdeaF72ms4MXzGsolLq';
  dept_it UUID;
  dept_fin UUID;
  div_itops UUID;
  div_itrnd UUID;
  div_hrd UUID;
  t_devops UUID;
  t_ppqa UUID;
  t_psops UUID;
  t_be UUID;
  t_fe UUID;
  t_qa UUID;
  t_uiux UUID;
  t_hrd UUID;
BEGIN
  -- Get IDs dynamically
  SELECT id INTO dept_it FROM departments WHERE name = 'IT' LIMIT 1;
  SELECT id INTO dept_fin FROM departments WHERE name = 'Finance' LIMIT 1;
  SELECT id INTO div_itops FROM divisions WHERE name = 'IT OPS' LIMIT 1;
  SELECT id INTO div_itrnd FROM divisions WHERE name = 'IT RnD' LIMIT 1;
  SELECT id INTO div_hrd FROM divisions WHERE name = 'HRD' LIMIT 1;
  SELECT id INTO t_devops FROM teams WHERE name = 'DevOps' LIMIT 1;
  SELECT id INTO t_ppqa FROM teams WHERE name = 'PPQA' LIMIT 1;
  SELECT id INTO t_psops FROM teams WHERE name = 'PSOps' LIMIT 1;
  SELECT id INTO t_be FROM teams WHERE name = 'Developer Backend' LIMIT 1;
  SELECT id INTO t_fe FROM teams WHERE name = 'Developer Frontend' LIMIT 1;
  SELECT id INTO t_qa FROM teams WHERE name = 'QA Tester' LIMIT 1;
  SELECT id INTO t_uiux FROM teams WHERE name = 'UI/UX' LIMIT 1;
  SELECT id INTO t_hrd FROM teams WHERE name = 'HRD' LIMIT 1;

  -- ============ MANAGERS (Approvers per Division) ============
  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Ryan IT-OPS Manager','ryan.manager@itsm.local',pw,'approver',true,dept_it,div_itops,NULL,'manager')
  ON CONFLICT(email) DO NOTHING;

  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Yudhit IT-RnD Manager','yudhit.manager@itsm.local',pw,'approver',true,dept_it,div_itrnd,NULL,'manager')
  ON CONFLICT(email) DO NOTHING;

  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Sari HRD Manager','sari.manager@itsm.local',pw,'approver',true,dept_fin,div_hrd,NULL,'manager')
  ON CONFLICT(email) DO NOTHING;

  -- ============ LEADERS (Agents) ============
  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Andi DevOps Lead','andi.lead@itsm.local',pw,'agent',true,dept_it,div_itops,t_devops,'leader')
  ON CONFLICT(email) DO NOTHING;

  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Budi PPQA Lead','budi.lead@itsm.local',pw,'agent',true,dept_it,div_itops,t_ppqa,'leader')
  ON CONFLICT(email) DO NOTHING;

  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Citra PSOps Lead','citra.lead@itsm.local',pw,'agent',true,dept_it,div_itops,t_psops,'leader')
  ON CONFLICT(email) DO NOTHING;

  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Deni Backend Lead','deni.lead@itsm.local',pw,'agent',true,dept_it,div_itrnd,t_be,'leader')
  ON CONFLICT(email) DO NOTHING;

  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Eka Frontend Lead','eka.lead@itsm.local',pw,'agent',true,dept_it,div_itrnd,t_fe,'leader')
  ON CONFLICT(email) DO NOTHING;

  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Fani QA Lead','fani.lead@itsm.local',pw,'agent',true,dept_it,div_itrnd,t_qa,'leader')
  ON CONFLICT(email) DO NOTHING;

  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Gita UIUX Lead','gita.lead@itsm.local',pw,'agent',true,dept_it,div_itrnd,t_uiux,'leader')
  ON CONFLICT(email) DO NOTHING;

  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Hana HRD Lead','hana.lead@itsm.local',pw,'agent',true,dept_fin,div_hrd,t_hrd,'leader')
  ON CONFLICT(email) DO NOTHING;

  -- ============ STAFF (Agents) ============
  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Irfan DevOps Staff','irfan.staff@itsm.local',pw,'agent',true,dept_it,div_itops,t_devops,'staff')
  ON CONFLICT(email) DO NOTHING;

  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Joko PPQA Staff','joko.staff@itsm.local',pw,'agent',true,dept_it,div_itops,t_ppqa,'staff')
  ON CONFLICT(email) DO NOTHING;

  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Kiki PSOps Staff','kiki.staff@itsm.local',pw,'agent',true,dept_it,div_itops,t_psops,'staff')
  ON CONFLICT(email) DO NOTHING;

  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Lina Backend Staff','lina.staff@itsm.local',pw,'agent',true,dept_it,div_itrnd,t_be,'staff')
  ON CONFLICT(email) DO NOTHING;

  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Mira Frontend Staff','mira.staff@itsm.local',pw,'agent',true,dept_it,div_itrnd,t_fe,'staff')
  ON CONFLICT(email) DO NOTHING;

  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Nanda QA Staff','nanda.staff@itsm.local',pw,'agent',true,dept_it,div_itrnd,t_qa,'staff')
  ON CONFLICT(email) DO NOTHING;

  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Omar UIUX Staff','omar.staff@itsm.local',pw,'agent',true,dept_it,div_itrnd,t_uiux,'staff')
  ON CONFLICT(email) DO NOTHING;

  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'Putri HRD Staff','putri.staff@itsm.local',pw,'agent',true,dept_fin,div_hrd,t_hrd,'staff')
  ON CONFLICT(email) DO NOTHING;

  -- ============ REGULAR USERS ============
  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'User IT','user.it@itsm.local',pw,'user',true,dept_it,NULL,NULL,NULL)
  ON CONFLICT(email) DO NOTHING;

  INSERT INTO users (id,full_name,email,password,role,is_active,department_id,division_id,team_id,position)
  VALUES (gen_random_uuid(),'User Finance','user.finance@itsm.local',pw,'user',true,dept_fin,NULL,NULL,NULL)
  ON CONFLICT(email) DO NOTHING;

END $$;
