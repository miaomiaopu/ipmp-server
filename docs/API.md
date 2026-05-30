# IPMP API 设计文档

## 概述

- Base URL: `/api/v1`
- 安全约束: **所有接口仅允许 GET 和 POST 方法**
- 统一响应格式:
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {},
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  }
  ```
- 认证: `Authorization: Bearer <access_token>`

## 错误码

| code | 含义 |
|------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或 Token 过期 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 405 | 不允许的 HTTP 方法 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |

---

## 认证 (Auth)

### POST `/auth/login`

登录获取 Token。登录限流：每 IP 每分钟最多 5 次。

```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "your_password"
}

Response 200:
{
  "code": 0,
  "data": {
    "access_token": "eyJhbG...",
    "refresh_token": "eyJhbG...",
    "expires_in": 86400,
    "user": {
      "id": "uuid",
      "username": "admin",
      "display_name": "管理员",
      "role": "admin"
    }
  }
}
```

### POST `/auth/refresh`

刷新 Access Token。

```
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbG..."
}
```

### POST `/auth/logout`

登出，注销 Refresh Token。

### GET `/auth/me`

获取当前用户信息。

---

## 用户 (Users)

### GET `/users`

用户列表（分页+搜索）。需要 admin 角色。

### GET `/users/:id`

用户详情。

### POST `/users`

创建用户（admin）。

```json
{
  "username": "zhangsan",
  "password": "secure_password",
  "display_name": "张三",
  "email": "zhangsan@example.com",
  "phone": "13800138000",
  "role": "user"
}
```

### POST `/users/:id/update`

更新用户信息。

### POST `/users/:id/delete`

软删除用户（admin）。

---

## 用户 AI 配置 (User AI Config)

每个用户独立管理自己的 AI Key。API Key 加密存储，**API 响应永远不返回完整 Key**。

### GET `/users/:id/ai-config`

获取当前 AI 配置。仅返回掩码后的 Key。

```json
{
  "code": 0,
  "data": {
    "configured": true,
    "provider": "deepseek",
    "model": "deepseek-v4-flash",
    "key_preview": "sk-****abc",
    "is_active": true
  }
}
```

### POST `/users/:id/ai-config/update`

更新 AI 配置。`api_key` 明文传输（通过 HTTPS 加密），服务端 AES-256-GCM 加密存储。

```json
{
  "provider": "deepseek",
  "api_key": "sk-xxxx",
  "model": "deepseek-v4-flash",
  "base_url": ""
}
```

### POST `/users/:id/ai-config/delete`

删除 AI 配置。

### POST `/users/:id/ai-config/test`

测试 AI 连接是否正常。

---

## 客户 (Customers)

### GET `/customers`

客户列表（分页+搜索）。查询参数：`?keyword=xxx&status=active&page=1&page_size=20`。

### GET `/customers/:id`

客户详情，含关联项目数量统计。

### POST `/customers`

创建客户。

```json
{
  "customer_code": "CUST-001",
  "name": "某某科技有限公司",
  "contact_person": "李四",
  "contact_phone": "13900139000",
  "contact_email": "lisi@example.com",
  "address": "北京市朝阳区xxx",
  "notes": "重要客户"
}
```

### POST `/customers/:id/update`

更新客户信息。

### POST `/customers/:id/delete`

软删除客户。

---

## 项目 (Projects)

### GET `/projects`

项目列表（分页+筛选）。查询参数：`?customer_id=xxx&status=in_progress&page=1&page_size=20`。

### GET `/projects/:id`

项目详情，含关联任务和需求数量。

### POST `/projects`

创建项目。

```json
{
  "project_code": "PROJ-001",
  "name": "ERP 管理系统",
  "customer_id": "uuid",
  "manager_id": "uuid",
  "start_date": "2026-05-01",
  "go_live_date": "2026-08-01",
  "completion_date": "",
  "status": "planning",
  "description": "企业资源管理系统开发"
}
```

### POST `/projects/:id/update`

更新项目信息。

### POST `/projects/:id/delete`

软删除项目。

---

## 任务 (Tasks)

统一任务模型，通过 `task_type` 区分类型。

### GET `/tasks`

任务列表（分页+筛选）。查询参数：`?task_type=project&project_id=xxx&status=todo&page=1&page_size=20`。

### GET `/tasks/:id`

任务详情，含工时汇总。

### POST `/tasks`

创建任务。

```json
{
  "task_type": "project",
  "title": "用户管理模块开发",
  "description": "实现用户CRUD功能",
  "project_id": "uuid",
  "customer_id": null,
  "status": "in_progress",
  "priority": "high",
  "due_date": "2026-06-15"
}
```

### POST `/tasks/:id/update`

更新任务信息。

### POST `/tasks/:id/status`

仅更新任务状态。

```json
{
  "status": "in_progress"
}
```

### POST `/tasks/:id/delete`

软删除任务。

---

## 需求 (Requirements)

### GET `/requirements`

需求列表（分页+筛选）。查询参数：`?req_type=project&project_id=xxx&status=approved`。

### GET `/requirements/:id`

需求详情。

### POST `/requirements`

创建需求。

```json
{
  "req_type": "project",
  "requirement_code": "0001",
  "title": "支持批量导入用户",
  "description": "通过 Excel 批量导入用户数据",
  "project_id": "uuid",
  "customer_id": null,
  "priority": "medium",
  "scheduled_date": "2026-06-20"
}
```

### POST `/requirements/:id/update`

更新需求。

### POST `/requirements/:id/status`

更新需求状态。

### POST `/requirements/:id/delete`

软删除需求。

---

## 工时记录 (Work Logs)

### GET `/work-logs`

工时列表（分页+筛选）。查询参数：`?user_id=xxx&project_id=xxx&start_date=2026-05-19&end_date=2026-05-25`。

### GET `/work-logs/:id`

工时详情。

### POST `/work-logs`

录入工时。

```json
{
  "task_id": "uuid",
  "log_date": "2026-05-27",
  "hours": 3.5,
  "description": "完成用户列表页面接口开发"
}
```

### POST `/work-logs/:id/update`

更新工时记录。

### POST `/work-logs/:id/delete`

删除工时记录。

### GET `/work-logs/stats`

工时聚合统计（核心报表接口）。查询参数：`?start_date=2026-05-19&end_date=2026-05-25&group_by=project`。

`group_by` 可选值：`day` | `week` | `project` | `customer`。

```json
{
  "code": 0,
  "data": [
    {
      "project_id": "uuid",
      "project_code": "PROJ-001",
      "project_name": "ERP 管理系统",
      "customer_code": "CUST-001",
      "customer_name": "某某科技有限公司",
      "total_hours": 40.0,
      "working_days": 5,
      "daily_avg_hours": 8.0,
      "entries": [
        {"date": "2026-05-25", "hours": 8.0, "description": "..."}
      ]
    }
  ]
}
```

### GET `/work-logs/export`

导出 Excel。查询参数：`?start_date=xxx&end_date=xxx&format=xlsx`。

---

## 周报 (Weekly Reports)

### GET `/weekly-reports`

周报列表。查询参数：`?user_id=xxx&week_start=2026-05-19&report_type=personal`。

### GET `/weekly-reports/:id`

周报详情。

### POST `/weekly-reports/generate`

生成周报（模板或 AI）。

```json
{
  "week_start": "2026-05-19",
  "week_end": "2026-05-25",
  "report_type": "personal",
  "project_id": null
}
```

### POST `/weekly-reports/:id/update`

编辑周报内容。

### POST `/weekly-reports/:id/review`

审核通过。

### POST `/weekly-reports/:id/finalize`

发布为终版。

---

## 仪表盘 (Dashboard)

### GET `/dashboard/stats`

```json
{
  "code": 0,
  "data": {
    "active_projects": 5,
    "tasks_this_week": 12,
    "hours_this_week": 40.0,
    "pending_requirements": 3
  }
}
```

### GET `/dashboard/this-week`

本周详细数据，按项目分组。

---

## AI 接口

### POST `/ai/generate-report`

使用当前用户配置的 AI Key 生成周报内容。如用户未配置 AI，返回错误。

```json
{
  "week_start": "2026-05-19",
  "week_end": "2026-05-25",
  "report_type": "personal",
  "project_id": null
}

Response 200:
{
  "code": 0,
  "data": {
    "content": "## 本周工作总结\n\n### 一、项目进展\n..."
  }
}
```

### POST `/ai/summarize`

通用文本摘要。

```json
{
  "text": "需要摘要的长文本"
}
```
