# 快速开始

## 环境要求

- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- [Git](https://git-scm.com/)

本地开发（不使用 Docker）额外需要：
- [Go](https://go.dev/dl/) 1.21+
- [Node.js](https://nodejs.org/) 20+

## 克隆并启动

```bash
# 克隆项目
git clone https://github.com/your-org/fullstack-engineering-lab.git
cd fullstack-engineering-lab

# 复制环境变量
cp .env.example .env

# 初始化（可选）
make init

# 启动所有服务
make up
```

## 验证服务

启动后，验证所有服务正常运行：

| 服务 | 地址 | 预期结果 |
|------|------|----------|
| 前端体验台 | http://localhost:3000 | 科技感深色主题首页 |
| 后端 API | http://localhost:8080/api/v1/health | `{"code":0,"message":"success",...}` |
| JWT 案例 | http://localhost:3000/cases/jwt-auth | 交互式演示页面 |
| 文档站 | http://localhost:5174 | VitePress 文档 |

## 常用命令

```bash
make up        # 启动所有服务
make down      # 停止所有服务
make logs      # 查看日志
make restart   # 重启所有服务
make clean     # 停止并清理数据卷
make test      # 运行测试
```

## 本地开发

如果不想使用 Docker，可以本地开发：

```bash
# 仅用 Docker 启动 MySQL 和 Redis
docker compose -f deploy/docker-compose.yml up -d mysql redis

# 启动 Go 后端
cd apps/server
go run ./cmd/server/

# 另一个终端，启动 Vue 前端
cd apps/web
npm install
npm run dev
```

前端将在 http://localhost:5173 启动，支持热更新。
