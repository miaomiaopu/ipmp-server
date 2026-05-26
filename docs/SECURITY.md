# IPMP 安全设计文档

## 概述

IPMP 面向公网部署，安全设计贯穿数据存储、API 通信、认证授权、运维部署全链路。

## 数据加密

### 敏感字段加密存储

以下字段在数据库中使用 **AES-256-GCM** 加密存储：

| 表 | 字段 | 说明 |
|----|------|------|
| users | phone, email | 用户联系方式 |
| customers | contact_phone, contact_email, address | 客户联系人信息 |
| user_ai_configs | api_key | 用户 AI Key |

**密钥管理**：
- 加密密钥通过环境变量 `IPMP_ENCRYPTION_KEY` 注入（32 字节 hex）
- 密钥不在代码、配置文件、日志中出现
- 服务启动时检测密钥是否设置，未设置则拒绝启动
- 密钥变更需执行数据迁移（解密→重加密）

**实现方式**：
- GORM 自定义类型 `EncryptedField`
- `BeforeSave` Hook 自动加密
- `AfterFind` Hook 自动解密
- 业务代码层完全无感

### 密码哈希

用户密码使用 bcrypt (cost=12) 哈希，不可逆。

## API 安全

### 认证机制

- **Access Token**: JWT，24 小时过期
- **Refresh Token**: JWT，7 天过期，存储于数据库可撤销
- Token 负载：`{ user_id, username, role, exp, iat }`
- 所有业务接口需 `Authorization: Bearer <token>` 头

### 方法白名单

- 全局中间件强制仅允许 GET 和 POST 方法
- 非 GET/POST 请求直接返回 `405 Method Not Allowed`

### 限流策略

| 接口 | 限制 |
|------|------|
| 登录 | 每 IP 每分钟 5 次 |
| 全局限流 | 每 IP 每秒 100 请求 |
| AI 接口 | 每用户每分钟 10 次 |

超限返回 `429 Too Many Requests`。

### 请求体大小限制

最大 10MB，防止大文件攻击。

## 传输安全

### HTTPS

生产环境强制 TLS 1.2+，Nginx 反向代理处理 TLS 终止。

### 响应安全头

| 头部 | 值 |
|------|-----|
| `X-Content-Type-Options` | `nosniff` |
| `X-Frame-Options` | `DENY` |
| `X-XSS-Protection` | `1; mode=block` |
| `Content-Security-Policy` | `default-src 'self'` |
| `Referrer-Policy` | `strict-origin-when-cross-origin` |
| `Permissions-Policy` | `camera=(), microphone=(), geolocation=()` |
| `Strict-Transport-Security` | `max-age=31536000; includeSubDomains` |

### CORS

仅允许配置中指定的前端域名（不允许 `*`）。

## 数据脱敏

API 响应根据角色对敏感字段进行脱敏：

| 角色 | 本人手机 | 他人手机 | 客户联系人 |
|------|----------|----------|------------|
| admin | 完整 | 完整 | 完整 |
| manager | 完整 | 末4位 | 完整(自己负责的客户) |
| user | 完整 | `****` | 末4位 |

AI Key 在任何情况下都不返回完整值，仅返回 `sk-****xxxx` 格式的掩码。

## 日志安全

- **敏感字段脱敏**: 日志中间件自动将 `password`、`api_key`、`token`、`secret` 等字段替换为 `[REDACTED]`
- **审计日志**: 所有写操作（POST）记录到 `audit_logs` 表
  - 字段：user_id, action, entity_type, entity_id, old_value, new_value, ip_address, user_agent
- 访问日志不在磁盘持久化，仅输出到 stdout（由容器日志收集）

## 数据库安全

- 仅允许 localhost/内网 IP 连接
- 应用使用专用数据库用户（非 superuser）
- 最小权限：仅 DML 权限，无 DDL（迁移时单独授权）
- 连接使用 TLS（如数据库支持）

## 运维安全

### 备份

- 数据库每日定时备份
- 备份文件使用 GPG 加密存储
- 保留最近 30 天备份

### 环境变量

- 所有密钥/密码通过环境变量注入
- `.env` 文件加入 `.gitignore`，永不提交
- 提供 `.env.example` 模板（仅占位符）

### 依赖安全

- 定期 `go mod tidy` + 安全扫描
- Docker 基础镜像使用 `alpine` 精简版
- 非 root 用户运行容器

## 安全检查清单

部署前逐项确认：

- [ ] `IPMP_ENCRYPTION_KEY` 已设置为强随机值
- [ ] `JWT_SECRET` 已设置为 64+ 字符随机值
- [ ] `DB_PASSWORD` 非默认值
- [ ] 默认 admin 密码已修改
- [ ] HTTPS 证书已配置
- [ ] Nginx 安全头已生效（curl -I 验证）
- [ ] CORS 白名单限制为实际前端域名
- [ ] 防火墙仅开放 443 端口（80 重定向到 443）
- [ ] 数据库备份计划已启用
