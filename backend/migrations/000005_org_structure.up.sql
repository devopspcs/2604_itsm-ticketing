-- departments
CREATE TABLE departments (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL UNIQUE,
    code       VARCHAR(50)  NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- divisions
CREATE TABLE divisions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    department_id UUID NOT NULL REFERENCES departments(id) ON DELETE RESTRICT,
    name          VARCHAR(255) NOT NULL,
    code          VARCHAR(50)  NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (department_id, code)
);

-- teams
CREATE TABLE teams (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    division_id UUID NOT NULL REFERENCES divisions(id) ON DELETE RESTRICT,
    name        VARCHAR(255) NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Add org fields to users (nullable for backward compatibility)
ALTER TABLE users
    ADD COLUMN department_id UUID REFERENCES departments(id) ON DELETE RESTRICT,
    ADD COLUMN division_id   UUID REFERENCES divisions(id)   ON DELETE RESTRICT,
    ADD COLUMN team_id       UUID REFERENCES teams(id)       ON DELETE RESTRICT,
    ADD COLUMN position      VARCHAR(50) CHECK (position IN ('division_manager','manager','leader','staff'));

CREATE INDEX idx_users_department_id ON users(department_id);
CREATE INDEX idx_users_division_id   ON users(division_id);
CREATE INDEX idx_users_team_id       ON users(team_id);
CREATE INDEX idx_users_position      ON users(position);
