-- MySQL 手动约束修复脚本。
-- 执行前必须先运行 check_constraints_mysql.sql，并确认没有活跃 code 重复。

ALTER TABLE customers DROP INDEX customer_code;
ALTER TABLE projects DROP INDEX project_code;

ALTER TABLE customers
    ADD COLUMN customer_code_alive VARCHAR(32)
        GENERATED ALWAYS AS (IF(deleted_at IS NULL, customer_code, NULL)) STORED,
    ADD UNIQUE KEY uk_customers_code_alive (customer_code_alive);

ALTER TABLE projects
    ADD COLUMN project_code_alive VARCHAR(32)
        GENERATED ALWAYS AS (IF(deleted_at IS NULL, project_code, NULL)) STORED,
    ADD UNIQUE KEY uk_projects_code_alive (project_code_alive);
