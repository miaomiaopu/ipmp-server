-- IPMP MySQL / TDSQL-C 初始化迁移
-- 9 张核心表 + 索引

-- 1. 用户表
CREATE TABLE users (
    id             CHAR(36)     PRIMARY KEY DEFAULT (UUID()),
    username       VARCHAR(64)  NOT NULL UNIQUE,
    password_hash  VARCHAR(256) NOT NULL,
    display_name   VARCHAR(128) NOT NULL DEFAULT '',
    email          TEXT         NOT NULL,
    phone          TEXT         NOT NULL,
    role           VARCHAR(32)  NOT NULL DEFAULT 'user',
    status         VARCHAR(32)  NOT NULL DEFAULT 'active',
    created_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at     DATETIME(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);

-- 2. 客户表
CREATE TABLE customers (
    id              CHAR(36)     PRIMARY KEY DEFAULT (UUID()),
    customer_code   VARCHAR(32)  NOT NULL UNIQUE,
    name            VARCHAR(256) NOT NULL,
    contact_person  VARCHAR(128) NOT NULL DEFAULT '',
    contact_phone   TEXT         NOT NULL,
    contact_email   TEXT         NOT NULL,
    address         TEXT         NOT NULL,
    notes           TEXT         NOT NULL,
    status          VARCHAR(32)  NOT NULL DEFAULT 'active',
    created_at      DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at      DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at      DATETIME(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_customers_code ON customers(customer_code);
CREATE INDEX idx_customers_status ON customers(status);
CREATE INDEX idx_customers_name ON customers(name);

-- 3. 项目表
CREATE TABLE projects (
    id              CHAR(36)     PRIMARY KEY DEFAULT (UUID()),
    project_code    VARCHAR(32)  NOT NULL UNIQUE,
    name            VARCHAR(256) NOT NULL,
    customer_id     CHAR(36),
    manager_id      CHAR(36),
    start_date      DATE,
    go_live_date    DATE,
    completion_date DATE,
    status          VARCHAR(32)  NOT NULL DEFAULT 'planning',
    description     TEXT         NOT NULL,
    created_at      DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at      DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at      DATETIME(3),
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (manager_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_projects_customer ON projects(customer_id);
CREATE INDEX idx_projects_manager ON projects(manager_id);
CREATE INDEX idx_projects_status ON projects(status);

-- 4. 任务表（统一任务模型）
CREATE TABLE tasks (
    id              CHAR(36)     PRIMARY KEY DEFAULT (UUID()),
    task_type       VARCHAR(32)  NOT NULL DEFAULT 'project',
    title           VARCHAR(256) NOT NULL,
    description     TEXT         NOT NULL,
    project_id      CHAR(36),
    customer_id     CHAR(36),
    status          VARCHAR(32)  NOT NULL DEFAULT 'in_progress',
    priority        VARCHAR(32)  NOT NULL DEFAULT 'medium',
    due_date        DATE,
    created_at      DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at      DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at      DATETIME(3),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_tasks_type ON tasks(task_type);
CREATE INDEX idx_tasks_project ON tasks(project_id);
CREATE INDEX idx_tasks_customer ON tasks(customer_id);
CREATE INDEX idx_tasks_status ON tasks(status);

-- 5. 需求表（统一需求模型）
CREATE TABLE requirements (
    id               CHAR(36)     PRIMARY KEY DEFAULT (UUID()),
    req_type         VARCHAR(32)  NOT NULL DEFAULT 'project',
    requirement_code VARCHAR(32)  NOT NULL DEFAULT '',
    title            VARCHAR(256) NOT NULL,
    description      TEXT         NOT NULL,
    project_id       CHAR(36),
    customer_id      CHAR(36),
    priority         VARCHAR(32)  NOT NULL DEFAULT 'medium',
    status           VARCHAR(32)  NOT NULL DEFAULT 'pending',
    scheduled_date   DATE,
    created_at       DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at       DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at       DATETIME(3),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_requirements_type ON requirements(req_type);
CREATE INDEX idx_requirements_project ON requirements(project_id);
CREATE INDEX idx_requirements_customer ON requirements(customer_id);

-- 6. 工时记录表
CREATE TABLE work_logs (
    id              CHAR(36)     PRIMARY KEY DEFAULT (UUID()),
    task_id         CHAR(36),
    user_id         CHAR(36)     NOT NULL,
    project_id      CHAR(36),
    customer_id     CHAR(36),
    log_date        DATE         NOT NULL,
    hours           DECIMAL(5,2) NOT NULL,
    description     TEXT         NOT NULL,
    created_at      DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at      DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at      DATETIME(3),
    FOREIGN KEY (task_id) REFERENCES tasks(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_work_logs_user ON work_logs(user_id);
CREATE INDEX idx_work_logs_task ON work_logs(task_id);
CREATE INDEX idx_work_logs_project ON work_logs(project_id);
CREATE INDEX idx_work_logs_customer ON work_logs(customer_id);
CREATE INDEX idx_work_logs_date ON work_logs(log_date);
CREATE INDEX idx_work_logs_user_date ON work_logs(user_id, log_date);

-- 7. 周报表
CREATE TABLE weekly_reports (
    id              CHAR(36)     PRIMARY KEY DEFAULT (UUID()),
    user_id         CHAR(36)     NOT NULL,
    week_start      DATE         NOT NULL,
    week_end        DATE         NOT NULL,
    report_type     VARCHAR(32)  NOT NULL DEFAULT 'personal',
    project_id      CHAR(36),
    content         TEXT         NOT NULL,
    ai_raw_content  TEXT,
    status          VARCHAR(32)  NOT NULL DEFAULT 'draft',
    created_at      DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at      DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at      DATETIME(3),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (project_id) REFERENCES projects(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_weekly_reports_user ON weekly_reports(user_id);
CREATE INDEX idx_weekly_reports_week ON weekly_reports(week_start, week_end);
CREATE INDEX idx_weekly_reports_project ON weekly_reports(project_id);

-- 8. 用户 AI 配置表
CREATE TABLE user_ai_configs (
    id         CHAR(36)     PRIMARY KEY DEFAULT (UUID()),
    user_id    CHAR(36)     NOT NULL UNIQUE,
    provider   VARCHAR(32)  NOT NULL DEFAULT 'deepseek',
    api_key    TEXT         NOT NULL,
    model      VARCHAR(128) NOT NULL DEFAULT '',
    base_url   VARCHAR(512) NOT NULL DEFAULT '',
    is_active  TINYINT(1)   NOT NULL DEFAULT 1,
    created_at DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 9. 审计日志表
CREATE TABLE audit_logs (
    id           CHAR(36)     PRIMARY KEY DEFAULT (UUID()),
    user_id      CHAR(36),
    action       VARCHAR(64)  NOT NULL,
    resource     VARCHAR(64)  NOT NULL,
    resource_id  VARCHAR(64),
    detail       JSON,
    ip_address   VARCHAR(64)  NOT NULL DEFAULT '',
    created_at   DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at);

-- 插入默认 admin 用户 (密码: admin123, bcrypt cost=12)
INSERT INTO users (id, username, password_hash, display_name, role, status)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'admin',
    '$2a$12$LJ3m4ys3GZfnYGecFX0rEOI1GvwYfe6mBLsBcEar3cGF8dQs2oPXi',
    '管理员',
    'admin',
    'active'
);
