# IPMP 数据模型文档 v0.2.0

## 通用基础字段
所有表继承 BaseModel: `id`(UUID) / `created_at` / `updated_at` / `deleted_at`(软删除)

## 模型字段清单

### users

| 字段 | 类型 | 变更 |
|------|------|------|
| `username` | VARCHAR(64) UNIQUE NOT NULL | |
| `password_hash` | VARCHAR(256) NOT NULL, JSON隐藏 | |
| `display_name` | VARCHAR(128) | |
| `email` | TEXT, EncryptedField AES加密 | |
| `phone` | TEXT, EncryptedField AES加密 | |
| `role` | VARCHAR(32) DEFAULT 'user' | v0.2: 去manager |
| `status` | VARCHAR(32) DEFAULT 'active' | |

### customers

| 字段 | 类型 | 变更 |
|------|------|------|
| `customer_code` | VARCHAR(32) UNIQUE NOT NULL | |
| `name` | VARCHAR(256) NOT NULL | |
| `contact_person` | VARCHAR(128) | 展示脱敏 |
| `contact_phone` | TEXT, EncryptedField | |
| `contact_email` | TEXT, EncryptedField | |
| `address` | TEXT, EncryptedField | |
| `notes` | TEXT | |
| `status` | VARCHAR(32) DEFAULT 'active' | |

### projects

| 字段 | 类型 | 变更 |
|------|------|------|
| `project_code` | VARCHAR(32) UNIQUE NOT NULL | |
| `name` | VARCHAR(256) NOT NULL | |
| `customer_id` | CHAR(36) FK→customers | v0.2: **必填** |
| `manager_id` | CHAR(36) FK→users | |
| `start_date` | DATE | |
| `go_live_date` | DATE | |
| `completion_date` | DATE | |
| `status` | VARCHAR(32) DEFAULT 'planning' | v0.2: 增 `online` |
| `description` | TEXT | |

### tasks

| 字段 | 类型 | 变更 |
|------|------|------|
| `task_type` | VARCHAR(32) DEFAULT 'project' | |
| `title` | VARCHAR(256) NOT NULL | |
| `description` | TEXT | |
| `project_id` | CHAR(36) FK→projects | project类型必填 |
| `customer_id` | CHAR(36) FK→customers | customer类型必填 |
| `status` | VARCHAR(32) DEFAULT 'in_progress' | v0.2: 去todo |
| `priority` | VARCHAR(32) DEFAULT 'medium' | |
| `due_date` | DATE | |
| ~~`assignee_id`~~ | — | v0.2: **删除** |
| ~~`estimated_hours`~~ | — | v0.2: **删除** |
| ~~`actual_hours`~~ | — | v0.2: **删除** |

### requirements

| 字段 | 类型 | 变更 |
|------|------|------|
| `req_type` | VARCHAR(32) DEFAULT 'project' | |
| `requirement_code` | VARCHAR(32) | v0.2: **新增** |
| `title` | VARCHAR(256) NOT NULL | |
| `description` | TEXT | |
| `project_id` | CHAR(36) FK→projects | project类型必填 |
| `customer_id` | CHAR(36) FK→customers | after_sales类型必填 |
| `priority` | VARCHAR(32) DEFAULT 'medium' | |
| `status` | VARCHAR(32) DEFAULT 'pending' | |
| `scheduled_date` | DATE | v0.2: **新增** |
| ~~`submitter`~~ | — | v0.2: **删除** |

---

## 迁移文件检查

| 表 | 迁移(MySQL) | 模型(GORM) | 差异 |
|------|:--:|:--:|------|
| users | 6 | 6 | ✅ |
| customers | 7 | 7 | ✅ |
| projects | 8 | 8 | ✅ |
| tasks | 10 | 8 | ❌ 迁移多 assignee_id/estimated_hours/actual_hours |
| requirements | 7 | 8 | ❌ 迁移多 submitter / 缺 requirement_code/scheduled_date |
| work_logs | 5 | 5 | ✅ |
| weekly_reports | 6 | 6 | ✅ |
| user_ai_configs | 5 | 5 | ✅ |
| audit_logs | 5 | 5 | ✅ |

**GORM AutoMigrate** 运行时自动处理: tasks 保留旧列(不删数据), requirements 新增并删除 submitter。迁移文件未同步是因本次在 fix 分支修改, 合并 main 后统一更新。

## 关联汇总

```
customers ──1:n──→ projects (customer_id)
customers ──1:n──→ tasks (customer_id, customer类型)
customers ──1:n──→ requirements (customer_id, after_sales类型)
customers ──1:n──→ work_logs (customer_id, 冗余)

projects ──1:n──→ tasks (project_id, project类型)
projects ──1:n──→ requirements (project_id, project类型)
projects ──1:n──→ work_logs (project_id, 冗余)

users ──1:n──→ projects (manager_id)
users ──1:n──→ work_logs (user_id)
users ──1:1──→ user_ai_configs (user_id, CASCADE)

tasks ──1:n──→ work_logs (task_id)
```

**删除级联**: 客户软删除→仅提示; 项目软删除→关联任务/需求同步软删除; 用户软删除→工时保留。

**软删除+去重**: 已删记录(code)允许新建同code, GORM deleted_at联合索引实现。
