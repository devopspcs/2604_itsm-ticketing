-- Seed departments
INSERT INTO departments (id, name, code) VALUES
    (gen_random_uuid(), 'URAI',    'URAI'),
    (gen_random_uuid(), 'Finance', 'FIN'),
    (gen_random_uuid(), 'SCO',     'SCO'),
    (gen_random_uuid(), 'KJS',     'KJS'),
    (gen_random_uuid(), 'NOS',     'NOS'),
    (gen_random_uuid(), 'IT',      'IT')
ON CONFLICT (name) DO NOTHING;

-- Seed divisions under IT
INSERT INTO divisions (id, department_id, name, code) VALUES
    (gen_random_uuid(), (SELECT id FROM departments WHERE code='IT'), 'IT RnD', 'IT-RND'),
    (gen_random_uuid(), (SELECT id FROM departments WHERE code='IT'), 'IT OPS', 'IT-OPS')
ON CONFLICT (department_id, code) DO NOTHING;

-- Seed teams under IT RnD
INSERT INTO teams (id, division_id, name)
SELECT gen_random_uuid(), d.id, t.name
FROM divisions d, (VALUES ('UI/UX'), ('Developer Frontend'), ('Developer Backend'), ('QA Tester')) AS t(name)
WHERE d.code = 'IT-RND'
ON CONFLICT DO NOTHING;

-- Seed teams under IT OPS
INSERT INTO teams (id, division_id, name)
SELECT gen_random_uuid(), d.id, t.name
FROM divisions d, (VALUES ('DevOps'), ('PSOps'), ('PPQA')) AS t(name)
WHERE d.code = 'IT-OPS'
ON CONFLICT DO NOTHING;
