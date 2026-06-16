-- PostgreSQL 手动约束修复脚本。
-- 执行前必须先运行 check_constraints_pg.sql，并确认没有活跃 code 重复。

BEGIN;

ALTER TABLE customers DROP CONSTRAINT IF EXISTS customers_customer_code_key;
ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_project_code_key;

DROP INDEX IF EXISTS idx_customers_code;
DROP INDEX IF EXISTS idx_projects_code;

CREATE UNIQUE INDEX IF NOT EXISTS uk_customers_code_alive
    ON customers(customer_code)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uk_projects_code_alive
    ON projects(project_code)
    WHERE deleted_at IS NULL;

COMMIT;
