DELETE FROM teams WHERE division_id IN (SELECT id FROM divisions WHERE code IN ('IT-RND', 'IT-OPS'));
DELETE FROM divisions WHERE code IN ('IT-RND', 'IT-OPS');
DELETE FROM departments WHERE code IN ('URAI', 'FIN', 'SCO', 'KJS', 'NOS', 'IT');
