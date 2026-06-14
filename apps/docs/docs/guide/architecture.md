# 项目架构

## Monorepo 结构

```
fullstack-engineering-lab/
├── apps/
│   ├── web/          # Vue 3 + Vite + TypeScript 前端体验台
│   ├── server/       # Go + Gin 后端 API 服务
│   └── docs/         # VitePress 文档站
├── examples/         # 案例说明文档与扩展目录
├── deploy/           # Docker Compose、Nginx、MySQL、Redis 配置
└── scripts/          # 辅助脚本
```

## 后端架构

Go 后端采用分层架构：

```
cmd/server/main.go        # 程序入口
internal/
├── config/               # 配置管理（Viper）
├── logger/               # 结构化日志（Zap）
├── database/             # 数据库连接（GORM）
├── model/                # 数据模型与 DTO
├── repository/           # 数据访问层
├── service/              # 业务逻辑层
├── handler/              # HTTP 处理器
├── router/               # 路由定义
├── middleware/            # JWT 认证、CORS 中间件
└── response/             # 统一响应格式
pkg/
├── jwt/                  # JWT 工具
└── password/             # bcrypt 密码工具
```

### 依赖注入链

```
配置 → 日志 → 数据库 → Redis → Repository → Service → Handler → Router → 启动
```

每一层通过构造函数接收依赖，保证代码可测试性。

## 前端架构

```
src/
├── api/              # Axios 实例与 API 函数
├── components/       # 可复用 UI 组件
├── layouts/          # 页面布局
├── router/           # Vue Router 路由配置
├── stores/           # Pinia 状态管理
├── styles/           # CSS 变量、全局样式、主题
├── types/            # TypeScript 类型定义
├── utils/            # 工具函数
└── views/            # 页面组件
```

## 数据流

```
浏览器 → Nginx → /api/* → Go 后端 → MySQL
                 → /*    → Vue 前端
                 → /docs → VitePress

客户端 → Axios → Bearer Token → JWT 中间件 → Handler → Service → Repository → 数据库
```

## 服务架构（Docker）

```
                    ┌──────────┐
                    │  Nginx   │ :80
                    └────┬─────┘
                         │
          ┌──────────────┼──────────────┐
          │              │              │
    ┌─────┴─────┐  ┌─────┴─────┐  ┌────┴────┐
    │  Vue 前端 │  │  Go 后端  │  │ 文档站  │
    │   :3000   │  │   :8080   │  │  :5174  │
    └───────────┘  └─────┬─────┘  └─────────┘
                         │
                  ┌──────┴──────┐
                  │             │
            ┌─────┴─────┐ ┌────┴────┐
            │   MySQL   │ │  Redis  │
            │   :3306   │ │  :6379  │
            └───────────┘ └─────────┘
```
