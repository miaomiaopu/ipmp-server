-- IPMP PostgreSQL 初始化迁移
-- 9 张核心表 + 索引

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. 用户表
CREATE TABLE users (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username       VARCHAR(64)  NOT NULL UNIQUE,
    password_hash  VARCHAR(256) NOT NULL,
    display_name   VARCHAR(128) NOT NULL DEFAULT '',
    email          TEXT         NOT NULL DEFAULT '',
    phone          TEXT         NOT NULL DEFAULT '',
    role           VARCHAR(32)  NOT NULL DEFAULT 'user',
    status         VARCHAR(32)  NOT NULL DEFAULT 'active',
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ
);
CREATE INDEX idx_users_username ON users(username) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;

-- 2. 客户表
CREATE TABLE customers (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_code   VARCHAR(32)  NOT NULL UNIQUE,
    name            VARCHAR(256) NOT NULL,
    contact_person  VARCHAR(128) NOT NULL DEFAULT '',
    contact_phone   TEXT         NOT NULL DEFAULT '',
    contact_email   TEXT         NOT NULL DEFAULT '',
    address         TEXT         NOT NULL DEFAULT '',
    notes           TEXT         NOT NULL DEFAULT '',
    status          VARCHAR(32)  NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX idx_customers_code ON customers(customer_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_customers_status ON customers(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_customers_name ON customers(name) WHERE deleted_at IS NULL;

-- 3. 项目表
CREATE TABLE projects (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_code    VARCHAR(32)  NOT NULL UNIQUE,
    name            VARCHAR(256) NOT NULL,
    customer_id     UUID         REFERENCES customers(id),
    manager_id      UUID         ,
    start_date      DATE,
    go_live_date    DATE,
    completion_date DATE,
    status          VARCHAR(32)  NOT NULL DEFAULT 'planning',
    description     TEXT         NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX idx_projects_customer ON projects(customer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_projects_manager ON projects(manager_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_projects_status ON projects(status) WHERE deleted_at IS NULL;

-- 4. 任务表（统一任务模型）
CREATE TABLE tasks (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_type       VARCHAR(32)  NOT NULL DEFAULT 'project',
    title           VARCHAR(256) NOT NULL,
    description     TEXT         NOT NULL DEFAULT '',
    project_id      UUID         REFERENCES projects(id),
    customer_id     UUID         REFERENCES customers(id),
    status          VARCHAR(32)  NOT NULL DEFAULT 'in_progress',
    priority        VARCHAR(32)  NOT NULL DEFAULT 'medium',
    due_date        DATE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX idx_tasks_type ON tasks(task_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_tasks_project ON tasks(project_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tasks_customer ON tasks(customer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tasks_status ON tasks(status) WHERE deleted_at IS NULL;

-- 5. 需求表（统一需求模型）
CREATE TABLE requirements (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    req_type         VARCHAR(32)  NOT NULL DEFAULT 'project',
    requirement_code VARCHAR(32)  NOT NULL DEFAULT '',
    title            VARCHAR(256) NOT NULL,
    description      TEXT         NOT NULL DEFAULT '',
    project_id       UUID         REFERENCES projects(id),
    customer_id      UUID         REFERENCES customers(id),
    priority         VARCHAR(32)  NOT NULL DEFAULT 'medium',
    status           VARCHAR(32)  NOT NULL DEFAULT 'pending',
    scheduled_date   DATE,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX idx_requirements_type ON requirements(req_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_requirements_project ON requirements(project_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_requirements_customer ON requirements(customer_id) WHERE deleted_at IS NULL;

-- 6. 工时记录表
CREATE TABLE work_logs (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id         UUID         REFERENCES tasks(id),
    user_id         UUID         NOT NULL ,
    project_id      UUID         REFERENCES projects(id),
    customer_id     UUID         REFERENCES customers(id),
    log_date        DATE         NOT NULL,
    hours           DECIMAL(5,2) NOT NULL,
    description     TEXT         NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX idx_work_logs_user ON work_logs(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_work_logs_task ON work_logs(task_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_work_logs_project ON work_logs(project_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_work_logs_customer ON work_logs(customer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_work_logs_date ON work_logs(log_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_work_logs_user_date ON work_logs(user_id, log_date) WHERE deleted_at IS NULL;

-- 7. 周报表
CREATE TABLE weekly_reports (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID         NOT NULL ,
    week_start      DATE         NOT NULL,
    week_end        DATE         NOT NULL,
    report_type     VARCHAR(32)  NOT NULL DEFAULT 'personal',
    project_id      UUID         REFERENCES projects(id),
    content         TEXT         NOT NULL DEFAULT '',
    ai_raw_content  TEXT,
    status          VARCHAR(32)  NOT NULL DEFAULT 'draft',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX idx_weekly_reports_user ON weekly_reports(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_weekly_reports_week ON weekly_reports(week_start, week_end) WHERE deleted_at IS NULL;
CREATE INDEX idx_weekly_reports_project ON weekly_reports(project_id) WHERE deleted_at IS NULL;

-- 8. 用户 AI 配置表
CREATE TABLE user_ai_configs (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID         NOT NULL UNIQUE  ON DELETE CASCADE,
    provider   VARCHAR(32)  NOT NULL DEFAULT 'deepseek',
    api_key    TEXT         NOT NULL,
    model      VARCHAR(128) NOT NULL DEFAULT '',
    base_url   VARCHAR(512) NOT NULL DEFAULT '',
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- 9. 审计日志表
CREATE TABLE audit_logs (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      UUID,
    action       VARCHAR(64)  NOT NULL,
    resource     VARCHAR(64)  NOT NULL,
    resource_id  VARCHAR(64),
    detail       JSONB,
    ip_address   VARCHAR(64)  NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at);

-- 插入默认 admin 用户 (密码: admin123, bcrypt cost=12)
-- 此密码哈希仅为初始化使用, 生产环境请立即修改
INSERT INTO users (id, username, password_hash, display_name, role, status)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'admin',
    '$2a$12$LJ3m4ys3GZfnYGecFX0rEOI1GvwYfe6mBLsBcEar3cGF8dQs2oPXi',
    '管理员',
    'admin',
    'active'
);
