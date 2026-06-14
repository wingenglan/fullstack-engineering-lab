# 贡献指南

欢迎贡献新的案例！以下是参与方式。

## 快速开始

1. Fork 本仓库
2. 克隆你的 Fork
3. 创建特性分支
4. 修改代码
5. 提交 Pull Request

## 开发环境搭建

```bash
# 克隆你的 Fork
git clone https://github.com/YOUR_USERNAME/fullstack-engineering-lab.git
cd fullstack-engineering-lab

# 初始化
cp .env.example .env
make init

# 启动开发环境
make dev
```

## 添加新案例

添加新工程案例的步骤：

1. 在 `examples/your-case/` 下创建文档
2. 在 `apps/server/internal/` 实现后端 API
3. 在 `apps/web/src/views/cases/` 实现前端页面
4. 在 `apps/web/src/router/index.ts` 添加路由
5. 更新 `CasesView.vue` 和 `HomeView.vue` 中的案例列表
6. 在 `apps/docs/docs/cases/` 添加文档

## 代码规范

### Go

- 遵循 [Effective Go](https://go.dev/doc/effective_go) 规范
- 使用 `gofmt` 格式化代码
- 提交前运行 `go vet`
- 为新功能编写测试

### TypeScript / Vue

- 所有新代码使用 TypeScript
- 遵循 Vue 3 Composition API 模式
- 使用 UnoCSS 工具类
- 保持组件专注和可复用

### 提交规范

使用约定式提交：

```
feat: 添加新案例
fix: 修复认证中间件 bug
docs: 更新 JWT 案例文档
chore: 更新依赖
```

## Pull Request 规范

- 每个 PR 只包含一个功能
- 同步更新相关文档
- 确保 CI 通过
- 编写清晰的 PR 描述

## 有问题？

提交 Issue 讨论后再开始大规模修改。
