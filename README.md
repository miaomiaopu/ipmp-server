# IPMP Server — 项目管理系统后端

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)
![Gin](https://img.shields.io/badge/Gin-1.x-0099FF?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql)
![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql)
![CI](https://img.shields.io/badge/CI-PG%20%7C%20MySQL%20矩阵-green?logo=githubactions)
![Version](https://img.shields.io/badge/Version-0.2.0-blue)
![License](https://img.shields.io/badge/License-Apache%202.0-blue)

IPMP（Intelligent Project Management Platform）后端服务，提供项目管理、工时统计、周报生成等 API 服务。支持 PostgreSQL / MySQL 双数据库，一行配置切换。AI 周报优先支持 DeepSeek。

## 目录

- [特性](#特性)
- [技术栈](#技术栈)
- [快速开始](#快速开始)
- [数据库切换](#数据库切换)
- [项目结构](#项目结构)
- [API 概览](#api-概览)
- [架构设计](#架构设计)
- [CI/CD](#cicd)
- [Docker 部署](#docker-部署)
- [安全](#安全)
- [扩展计划](#扩展计划)
- [贡献指南](#贡献指南)
- [许可证](#许可证)

## 特性

- **多数据库支持**: PostgreSQL / MySQL 一行配置切换，GORM 透明适配
- **客户管理**: 录入和维护客户信息，敏感信息加密存储
- **项目管理**: 项目全生命周期管理（立项→实施→竣工）
- **统一任务模型**: 项目任务 / 客户任务 / 日常任务，一套模型三类场景
- **需求跟踪**: 项目需求 + 客户售后需求，状态流转管理
- **工时统计**: 按日录入工时，按周汇总，8 小时比例分配
- **周报生成**: 模板化生成个人/项目周报，支持 AI 生成（优先 DeepSeek）
- **用户独立 AI 配置**: 每个用户配置自己的 AI Key，加密存储，Key 不出服务器
- **数据安全**: AES-256-GCM 加密敏感字段，JWT 认证，审计日志

## 技术栈

| 类别 | 选型 |
|------|------|
| 语言 | Go 1.25 |
| Web 框架 | Gin |
| ORM | GORM (PG + MySQL 双驱动) |
| 数据库 | PostgreSQL 16 / MySQL 8.0 |
| 认证 | JWT (access + refresh token) |
| 加密 | AES-256-GCM |
| 配置 | Viper (`${ENV_VAR}` 格式) |
| 日志 | Zerolog |
| AI | DeepSeek (优先) / OpenAI / Claude (Provider 接口可扩展) |
| Excel | excelize v2 |

## 快速开始

### 环境要求

- Go 1.25+
- PostgreSQL 16+ 或 MySQL 8.0+

```bash
git clone <repo-url> ipmp-server
cd ipmp-server
cp .env.example .env    # 编辑 .env 填入数据库连接信息
go mod download
go run cmd/server/main.go
# 默认监听 http://localhost:8080
```

## 数据库切换

通过 `DB_TYPE` 环境变量切换数据库，无需修改代码：

| 变量 | PG 默认值 | MySQL 默认值 | 说明 |
|------|-----------|-------------|------|
| `DB_TYPE` | `postgres` | `mysql` | 切换开关 |
| `*_HOST` | `localhost` | `localhost` | 数据库地址 |
| `*_PORT` | `5432` | `3306` | 端口 |
| `*_USER` | `postgres` | `root` | 用户名 |
| `*_PASSWORD` | — | — | 密码 |
| `*_NAME` | `ipmp` | `ipmp` | 数据库名 |

```bash
# PostgreSQL (默认)
export DB_TYPE=postgres

# MySQL
export DB_TYPE=mysql
```

迁移文件分别存放于 `migrations/pg/` 和 `migrations/mysql/`，按 `DB_TYPE` 选择执行。

## 项目结构

```
ipmp-server/
├── cmd/server/main.go            # 应用入口
├── internal/
│   ├── config/config.go          # 配置管理 (Viper, 双DB)
│   ├── middleware/               # 安全中间件
│   ├── model/                    # GORM 数据模型
│   ├── dto/{request,response}/   # 请求/响应 DTO
│   ├── repository/               # 数据访问层
│   ├── service/                  # 业务逻辑层
│   │   └── ai/                   # AI Provider 接口
│   ├── handler/                  # HTTP 处理器
│   ├── router/router.go          # 路由注册
│   └── pkg/                      # 内部工具包
│       ├── crypto/               # AES 加密
│       └── jwt/                  # JWT 工具
├── migrations/
│   ├── pg/                       # PostgreSQL DDL
│   └── mysql/                    # MySQL DDL
├── .github/workflows/
│   ├── ci.yml                    # CI: PG + MySQL 双 job 测试
│   └── deploy-dev.yml            # CD-Dev: dev* tag 触发 → Docker 构建 → 部署
├── config.yaml
├── docker-compose.yml
├── Dockerfile
├── nginx.conf
└── .env.example
```

## API 概览

> 安全约束：所有接口仅允许 GET 和 POST 方法。

| 模块 | 端点示例 | 说明 |
|------|----------|------|
| Auth | `POST /api/v1/auth/login` | 登录认证 + 密码修改 |
| Users | `GET /api/v1/users` | 用户管理（admin 专属） |
| Customers | `GET /api/v1/customers` | 客户 CRUD |
| Projects | `GET /api/v1/projects` | 项目 CRUD |
| Tasks | `GET /api/v1/tasks` | 任务管理 + 状态流转 |
| Requirements | `GET /api/v1/requirements` | 需求管理 |
| WorkLogs | `GET /api/v1/work-logs/stats` | 工时录入 + 统计聚合 |
| AIConfig | `GET /api/v1/ai-config` | 用户独立 AI Key 配置 |
| Reports | `POST /api/v1/weekly-reports/generate` | 生成周报 |
| AI | `POST /api/v1/ai/generate-report` | AI 生成报告 (DeepSeek/OpenAI/Claude) |

详细 API 文档参见 [docs/API.md](docs/API.md)。

## 架构设计

```
Handler → Service → Repository → PostgreSQL / MySQL (DB_TYPE 切换)
              ↘ AI Provider (DeepSeek / OpenAI / Claude)
```

- **Handler**: 解析请求、调用 Service、返回响应，零业务逻辑
- **Service**: 业务规则、流程编排、AI 调用、数据聚合
- **Repository**: GORM 数据库操作，加密字段通过 GORM Hook 透明处理
- **Model**: 数据结构定义，UUID 主键，软删除，时间戳

## CI/CD

| Pipeline | 触发 | 内容 |
|----------|------|------|
| **CI** | push / PR to main | PG + MySQL 双 job 独立测试、golangci-lint、gitleaks、build |
| **CD-Dev** | tag `dev*` | Docker 构建 → ghcr.io → SSH 部署到开发服务器 |

## 分支策略

```
feat/xxx → PR → CI → merge main → tag dev* → CD-Dev 部署
```

## Dev 部署

详见 [docker-compose.yml](docker-compose.yml) 和 [nginx.conf](nginx.conf)。

```bash
# 服务器初始化（仅一次）
mkdir -p /opt/ipmp/www
cp .env.example .env   # 填入数据库/JWT/加密密钥

# 触发部署
git tag dev-0.1.0 && git push origin main --tags
# CD-Dev 自动: Docker 构建 → push ghcr.io → rsync config → SSH 部署
```

## 安全

- **传输安全**: 生产环境强制 HTTPS/TLS 1.2+
- **认证**: JWT 双 Token 机制，登录限流
- **加密存储**: 联系方式、AI Key 等敏感字段 AES-256-GCM 加密
- **API 脱敏**: 响应按角色掩码敏感信息
- **审计日志**: 记录所有写操作和敏感读取
- **安全响应头**: HSTS, CSP, X-Frame-Options 等
- **方法白名单**: 仅允许 GET/POST 方法

详见 [docs/SECURITY.md](docs/SECURITY.md)。

## 扩展计划

- [x] 客户/项目 CRUD + 多数据库 + CI/CD + Dev 部署 (Phase 1)
- [x] 用户管理 + 任务 + 需求 + 工时 + AI 配置 + admin 仪表盘 (Phase 2)
- [ ] 周报生成 + 工时统计 + Excel 导出 (Phase 3)
- [ ] AI 周报生成 — DeepSeek 优先 (Phase 4)
- [ ] 文件附件 + 通知系统 + 暗色模式 + 生产 CD (Phase 5)

## 贡献指南

1. Fork 本仓库
2. 从 `main` 创建 `feat/xxx` 分支
3. 确保通过 gitleaks 扫描 + CI 矩阵测试
4. 创建 Pull Request 到 main

## 许可证

本项目基于 Apache License 2.0 开源。详见 [LICENSE](LICENSE)。

Copyright 2026 miaomiaopu
