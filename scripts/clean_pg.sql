-- IPMP PostgreSQL 清库脚本（保留表结构 + admin 用户）
TRUNCATE TABLE audit_logs, user_ai_configs, work_logs, requirements, tasks, weekly_reports, projects, customers RESTART IDENTITY CASCADE;
DELETE FROM users WHERE username != 'admin';
