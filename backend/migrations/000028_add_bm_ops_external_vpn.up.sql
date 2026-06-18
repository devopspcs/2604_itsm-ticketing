INSERT INTO external_services (name, type, url, auth_username, auth_password, extra_config, is_active) VALUES
('VPN BM', 'pritunl', 'https://64.176.81.185', '', '', '{"org_id":"6660a0981e61de46388fe649"}', true),
('VPN Operation', 'pritunl', 'https://45.32.99.199', '', '', '{"org_id":"66cea989d9a53fff04bf42ad"}', true),
('VPN External', 'pritunl', 'https://149.28.128.72', '', '', '{"org_id":"66f2bdacac23c1315c34a167"}', true)
ON CONFLICT DO NOTHING;
