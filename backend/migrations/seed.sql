-- Run this after migrations to create default users
-- Password for all: Admin@123

-- Admin user
INSERT INTO users (id, full_name, email, password, role, is_active, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'System Admin',
    'admin@itsm.local',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj4J/HS.iK8i',
    'admin',
    true,
    NOW(),
    NOW()
) ON CONFLICT (email) DO NOTHING;

-- Approver user
INSERT INTO users (id, full_name, email, password, role, is_active, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'Default Approver',
    'approver@itsm.local',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj4J/HS.iK8i',
    'approver',
    true,
    NOW(),
    NOW()
) ON CONFLICT (email) DO NOTHING;

-- Regular user
INSERT INTO users (id, full_name, email, password, role, is_active, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'Default User',
    'user@itsm.local',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj4J/HS.iK8i',
    'user',
    true,
    NOW(),
    NOW()
) ON CONFLICT (email) DO NOTHING;
