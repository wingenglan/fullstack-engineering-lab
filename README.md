<p align="center">
  <h1 align="center">FullStack Engineering Lab</h1>
  <p align="center">可运行、可体验、可学习的全栈工程实践案例库</p>
  <p align="center">
    <img src="https://img.shields.io/badge/Vue-3-4FC08D?style=flat-square&logo=vue.js" />
    <img src="https://img.shields.io/badge/Go-1.21-00ADD8?style=flat-square&logo=go" />
    <img src="https://img.shields.io/badge/TypeScript-5-3178C6?style=flat-square&logo=typescript" />
    <img src="https://img.shields.io/badge/Docker-Compose-2496ED?style=flat-square&logo=docker" />
    <img src="https://img.shields.io/badge/License-MIT-6366F1?style=flat-square" />
  </p>
</p>

---

## 项目介绍

FullStack Engineering Lab 是一个面向开发者的 **技术实验室 / Developer Playground**，用于沉淀开发过程中常见技术的完整实战案例。

**这不是一个后台管理系统**，而是一个可以运行、体验和学习的工程实践案例库。

### 案例规划

| 阶段 | 案例 | 状态 |
|------|------|------|
| Phase 1 | JWT 认证授权 | 已完成 |
| Phase 2 | WebSocket 实时通信、Redis 分布式锁 | 部分完成 |
| Phase 3 | 消息队列、定时任务 | 计划中 |
| Phase 4 | 文件上传、大文件分片、支付对接 | 计划中 |
| Phase 5 | 搜索引擎、AI 集成 | 计划中 |
| Phase 6 | Docker 部署、可观测性 | 计划中 |

## 技术栈

### 前端

- **Vue 3** + **Vite** + **TypeScript**
- **Element Plus** + **UnoCSS**
- **Pinia** 状态管理
- **Vue Router** 路由
- **Axios** HTTP 请求
- **lucide-vue-next** 图标

### 后端

- **Go** + **Gin** + **Gorm**
- **Viper** 配置管理
- **Zap** 结构化日志
- **JWT** 认证
- **Swaggo** API 文档

### 基础设施

- **MySQL** 主数据库
- **Redis** 缓存 / Token 黑名单
- **Nginx** 反向代理
- **Docker Compose** 容器编排
- **VitePress** 文档站

## 目录结构

```
fullstack-engineering-lab/
├── apps/
│   ├── web/                    # Vue 3 + Vite + TypeScript 前端
│   ├── server/                 # Go + Gin 后端服务
│   └── docs/                   # VitePress 文档站
├── examples/                   # 技术案例说明与扩展
│   ├── jwt-auth/               # JWT 认证案例文档
│   ├── websocket/              # (规划中)
│   ├── redis-lock/             # (规划中)
│   └── ...
├── deploy/                     # 部署配置
│   ├── docker-compose.yml
│   ├── nginx/
│   ├── mysql/
│   └── redis/
├── scripts/                    # 工具脚本
├── .github/workflows/          # CI/CD
├── Makefile
└── README.md
```

## 快速开始

### 前置要求

- Docker & Docker Compose
- Git

### 一键启动

```bash
# 克隆项目
git clone https://github.com/your-org/fullstack-engineering-lab.git
cd fullstack-engineering-lab

# 复制环境变量
cp .env.example .env

# 启动所有服务
make up
# 或
docker compose -f deploy/docker-compose.yml up -d
```

### 访问服务

| 服务 | 地址 |
|------|------|
| 前端体验台 | http://localhost:3000 |
| 后端 API | http://localhost:8080 |
| JWT Demo | http://localhost:3000/cases/jwt-auth |
| Redis Lock Demo | http://localhost:3000/cases/redis-lock |
| 健康检查 | http://localhost:8080/api/v1/health |
| 文档站 | http://localhost:5174 |
| API 文档 | http://localhost:8080/swagger/index.html |

### 本地开发

```bash
# 初始化项目
make init

# 启动本地开发服务（不使用 Docker）
make dev
```

### 常用命令

```bash
make up        # 启动所有服务
make down      # 停止所有服务
make logs      # 查看日志
make restart   # 重启所有服务
make clean     # 停止并清理数据卷
make test      # 运行测试
make build     # 构建所有镜像
```

## JWT Auth Demo

第一个完整案例，演示 JWT 认证授权的完整流程：

- 用户注册（bcrypt 密码加密）
- 用户登录（获取 access token）
- 获取用户信息（JWT 中间件验证）
- 退出登录（Redis Token 黑名单）
- Token 解码展示
- 请求日志实时查看

## Redis Lock Demo

第二个案例，演示基于 Redis 的分布式锁：

- 获取分布式锁（SET NX EX）
- 释放分布式锁（Lua 原子操作）
- 锁状态实时查询
- 并发争抢演示（多协程互斥）
- TTL 自动过期防死锁
- Owner Token 防误删

## Roadmap

- **Phase 1**: Auth & Basic Engineering - JWT 认证、基础工程化
- **Phase 2**: Realtime & Cache - WebSocket 实时通信、Redis 分布式锁
- **Phase 3**: MQ & Scheduler - 消息队列、定时任务
- **Phase 4**: File & Payment - 文件上传、支付对接
- **Phase 5**: Search & AI - 搜索引擎、AI 集成
- **Phase 6**: DevOps & Observability - 部署、监控、可观测性

## 贡献指南

欢迎贡献新的案例！请参考 [贡献指南](examples/jwt-auth/README.md)。

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/awesome-case`)
3. 提交变更 (`git commit -m 'Add awesome case'`)
4. 推送分支 (`git push origin feature/awesome-case`)
5. 创建 Pull Request

## License

[MIT](LICENSE)
