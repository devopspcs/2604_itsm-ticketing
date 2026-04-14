-- Seed default users, password: adminpcs
INSERT INTO users (id, full_name, email, password, role, is_active, created_at, updated_at)
VALUES (
    gen_random_uuid(), 'Admin PCS', 'admin@itsm.local',
    '$2a$12$SQbU9JeG5UlGmoBx3yUaZODfA00D2g8t28HdeaF72ms4MXzGsolLq',
    'admin', true, NOW(), NOW()
) ON CONFLICT (email) DO NOTHING;

INSERT INTO users (id, full_name, email, password, role, is_active, created_at, updated_at)
VALUES (
    gen_random_uuid(), 'Default Approver', 'approver@itsm.local',
    '$2a$12$SQbU9JeG5UlGmoBx3yUaZODfA00D2g8t28HdeaF72ms4MXzGsolLq',
    'approver', true, NOW(), NOW()
) ON CONFLICT (email) DO NOTHING;

INSERT INTO users (id, full_name, email, password, role, is_active, created_at, updated_at)
VALUES (
    gen_random_uuid(), 'Default User', 'user@itsm.local',
    '$2a$12$SQbU9JeG5UlGmoBx3yUaZODfA00D2g8t28HdeaF72ms4MXzGsolLq',
    'user', true, NOW(), NOW()
) ON CONFLICT (email) DO NOTHING;
