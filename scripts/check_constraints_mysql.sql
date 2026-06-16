-- MySQL 只读约束检查：发现仍存在的单列 UNIQUE 和活跃 code 重复。

SHOW INDEX FROM customers WHERE Non_unique = 0;
SHOW INDEX FROM projects WHERE Non_unique = 0;

SELECT customer_code, COUNT(*) AS active_count
FROM customers
WHERE deleted_at IS NULL
GROUP BY customer_code
HAVING COUNT(*) > 1;

SELECT project_code, COUNT(*) AS active_count
FROM projects
WHERE deleted_at IS NULL
GROUP BY project_code
HAVING COUNT(*) > 1;
