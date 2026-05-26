# IPMP Server — 项目管理系统后端

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)
![Gin](https://img.shields.io/badge/Gin-1.x-0099FF?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql)
![License](https://img.shields.io/badge/License-Apache%202.0-blue)

IPMP（Intelligent Project Management Platform）后端服务，提供项目管理、工时统计、周报生成等 API 服务。

## 目录

- [特性](#特性)
- [技术栈](#技术栈)
- [快速开始](#快速开始)
- [项目结构](#项目结构)
- [API 概览](#api-概览)
- [架构设计](#架构设计)
- [安全](#安全)
- [扩展计划](#扩展计划)
- [贡献指南](#贡献指南)
- [许可证](#许可证)

## 特性

- **客户管理**: 录入和维护客户信息，敏感信息加密存储
- **项目管理**: 项目全生命周期管理（立项→实施→竣工）
- **统一任务模型**: 项目任务 / 客户任务 / 日常任务，一套模型三类场景
- **需求跟踪**: 项目需求 + 客户售后需求，状态流转管理
- **工时统计**: 按日录入工时，按周汇总，8 小时比例分配
- **周报生成**: 模板化生成个人/项目周报，预留 AI 生成能力
- **数据安全**: AES-256-GCM 加密敏感字段，JWT 认证，审计日志

## 技术栈

| 类别 | 选型 |
|------|------|
| 语言 | Go 1.25 |
| Web 框架 | Gin |
| ORM | GORM |
| 数据库 | PostgreSQL 16 |
| 认证 | JWT (access + refresh token) |
| 加密 | AES-256-GCM |
| 配置 | Viper |
| 日志 | Zerolog |
| Excel | excelize v2 |

## 快速开始

### 环境要求

- Go 1.25+
- PostgreSQL 16+
- (可选) Docker & Docker Compose

### 本地开发

```bash
# 克隆仓库
git clone <repo-url> ipmp-server
cd ipmp-server

# 复制环境变量
cp .env.example .env
# 编辑 .env 填入实际的数据库密码和加密密钥

# 安装依赖
go mod download

# 运行数据库迁移
go run cmd/migrate/main.go

# 启动开发服务器
go run cmd/server/main.go
# 默认监听 http://localhost:8080
```

### Docker 部署

```bash
docker-compose up -d
```

## 项目结构

```
ipmp-server/
├── cmd/
│   └── server/main.go            # 应用入口
├── internal/
│   ├── config/config.go          # 配置管理 (Viper)
│   ├── middleware/               # 中间件 (auth, cors, ratelimit, ...)
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
├── migrations/                   # SQL 迁移文件
├── config.yaml                   # 默认配置
├── Dockerfile
├── docker-compose.yml
├── .env.example
└── Makefile
```

## API 概览

> 安全约束：所有接口仅允许 GET 和 POST 方法。

| 模块 | 端点示例 | 说明 |
|------|----------|------|
| Auth | `POST /api/v1/auth/login` | 登录认证 |
| Customers | `GET /api/v1/customers` | 客户列表 |
| Projects | `GET /api/v1/projects` | 项目列表 |
| Tasks | `GET /api/v1/tasks` | 任务管理 |
| WorkLogs | `GET /api/v1/work-logs/stats` | 工时聚合统计 |
| Reports | `POST /api/v1/weekly-reports/generate` | 生成周报 |
| AI | `POST /api/v1/ai/generate-report` | AI 生成报告 |

详细 API 文档参见 [API 设计文档](docs/API.md)。

## 架构设计

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Handler    │ ──▶ │   Service    │ ──▶ │  Repository  │
│ (HTTP 处理)  │ ◀── │ (业务逻辑)   │ ◀── │  (数据访问)  │
└──────────────┘     └──────┬───────┘     └──────┬───────┘
                            │                    │
                     ┌──────▼───────┐     ┌──────▼───────┐
                     │  AI Provider │     │  PostgreSQL   │
                     │  (OpenAI等)  │     │               │
                     └──────────────┘     └──────────────┘
```

- **Handler**: 解析请求、调用 Service、返回响应，零业务逻辑
- **Service**: 业务规则、流程编排、AI 调用、数据聚合
- **Repository**: GORM 数据库操作，加密字段通过 GORM Hook 透明处理
- **Model**: 数据结构定义，含软删除、UUID 主键、时间戳

## 安全

- **传输安全**: 生产环境强制 HTTPS/TLS 1.2+
- **认证**: JWT 双 Token 机制，登录限流
- **加密存储**: 联系方式、AI Key 等敏感字段 AES-256-GCM 加密
- **API 脱敏**: 响应按角色掩码敏感信息
- **审计日志**: 记录所有写操作和敏感读取
- **安全响应头**: HSTS, CSP, X-Frame-Options 等

详见 [安全设计文档](docs/SECURITY.md)。

## 扩展计划

- [x] 核心 CRUD (客户/项目/任务/需求)
- [x] 工时统计与报表
- [ ] AI 周报生成 (Phase 4)
- [ ] 文件附件上传
- [ ] 多用户权限 (RBAC)
- [ ] 通知系统
- [ ] 数据导入导出

## 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 确保代码通过 gitleaks 密钥扫描
4. 提交变更 (`git commit -m 'feat: add amazing feature'`)
5. 推送到分支 (`git push origin feature/amazing-feature`)
6. 创建 Pull Request

## 许可证

本项目基于 Apache License 2.0 开源。详见 [LICENSE](LICENSE) 文件。
