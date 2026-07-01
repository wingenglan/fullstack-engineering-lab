# JWT 认证授权案例

## 概述

JWT 认证授权案例是工程实验室的第一个完整案例，演示了基于 JSON Web Token 的完整认证流程。

## 工作原理

### 用户注册
1. 用户提交用户名、邮箱和密码
2. 服务端校验输入并检查唯一性
3. 使用 bcrypt 对密码进行哈希处理（cost=10）
4. 将用户信息存入 MySQL

### 用户登录
1. 用户提交用户名和密码
2. 服务端查找用户并用 bcrypt 验证密码
3. 生成包含用户 ID、签发时间和过期时间的 JWT Token
4. 返回 access_token 和过期时间

### 身份验证
1. 客户端将 Token 存储在 localStorage
2. 每个请求在 Header 中携带 `Authorization: Bearer <token>`
3. JWT 中间件验证 Token 签名和有效期
4. 中间件检查 Redis 黑名单（已撤销的 Token）
5. 用户 ID 被注入请求上下文

### 退出登录
1. 客户端发送退出请求
2. 服务端将当前 Token 加入 Redis 黑名单
3. TTL 设置为 Token 剩余有效期
4. 后续使用该 Token 的请求将被拒绝

## 接口文档

### POST /api/v1/auth/register

用户注册。

```json
// 请求
{
  "username": "demo",
  "email": "demo@example.com",
  "password": "123456"
}

// 响应
{
  "code": 0,
  "message": "success",
  "data": null
}
```

### POST /api/v1/auth/login

用户登录，获取 JWT Token。

```json
// 响应
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 7200
  }
}
```

### GET /api/v1/auth/profile

获取当前用户信息（需要认证）。

```json
// 响应
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "demo",
    "email": "demo@example.com",
    "nickname": "",
    "status": 1,
    "created_at": "2026-06-08T00:00:00Z"
  }
}
```

### POST /api/v1/auth/logout

撤销当前 Token（需要认证）。

## 数据库设计

```sql
CREATE TABLE users (
  id            BIGINT PRIMARY KEY AUTO_INCREMENT,
  username      VARCHAR(64)  NOT NULL UNIQUE,
  email         VARCHAR(128) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  nickname      VARCHAR(64)  DEFAULT '',
  status        TINYINT      NOT NULL DEFAULT 1,
  created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at    DATETIME     NULL
);
```

## 安全说明

- 密码使用 bcrypt 哈希存储（cost=10），绝不存储明文
- JWT 使用 HS256 签名算法，密钥可配置
- Token 过期时间可通过 `JWT_EXPIRE_MINUTES` 配置
- Redis 黑名单机制防止已撤销的 Token 被使用
- 开发环境配置 CORS；生产环境由 Nginx 统一代理
