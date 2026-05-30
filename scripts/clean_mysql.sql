-- IPMP MySQL 清库脚本（保留表结构 + admin 用户）
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE audit_logs;
TRUNCATE TABLE user_ai_configs;
TRUNCATE TABLE work_logs;
TRUNCATE TABLE requirements;
TRUNCATE TABLE tasks;
TRUNCATE TABLE weekly_reports;
TRUNCATE TABLE projects;
TRUNCATE TABLE customers;
DELETE FROM users WHERE username != 'admin';
SET FOREIGN_KEY_CHECKS = 1;
