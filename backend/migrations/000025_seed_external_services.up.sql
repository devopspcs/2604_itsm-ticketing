-- Seed VPN services
INSERT INTO external_services (name, type, url, auth_username, auth_password, extra_config, is_active) VALUES
('VPN DevOps', 'pritunl', 'https://10.184.20.11', 'devopspcs', '!*%bXMq|EoYJF0wh', '{"org_id":"65a090ed01046f906ef218d1"}', true),
('VPN IT', 'pritunl', 'https://10.184.20.13', 'devopspcs', '!*%bXMq|EoYJF0wh', '{"org_id":"659b7e2851948a075375d28a"}', true),
('VPN IT DevOps', 'pritunl', 'https://10.88.12.10', 'devopspcs', '!*%bXMq|EoYJF0wh', '{"org_id":"67ac71201b670c5fb706373d"}', true),
('VPN IT RnD', 'pritunl', 'https://10.184.20.14', 'devopspcs', '!*%bXMq|EoYJF0wh', '{"org_id":"65d564d839853ac36ddc4cf5"}', true)
ON CONFLICT DO NOTHING;

-- Update POSe with token
UPDATE external_services SET auth_token = '51Fr9WJ4nt6bWpvBVpw6piyHbm1VSAsVjA5vfMtczahHOjO7dcbMW2gOhU2Ayr3o' WHERE type = 'pose' AND auth_token = '';
