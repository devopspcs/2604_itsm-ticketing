-- Rollback: remove all seeded users
DELETE FROM refresh_tokens;
DELETE FROM users;
