# JWT Auth API 文档

## 基础信息

- Base URL: `/api/v1`
- Content-Type: `application/json`
- 认证方式: Bearer Token（部分接口）

## 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 错误码

| 错误码 | 含义 |
|--------|------|
| 0 | 成功 |
| 40001 | 认证失败（用户名或密码错误、无效 Token） |
| 40003 | Token 已过期 |
| 40004 | 无权限访问 |
| 50000 | 服务器内部错误 |

---

## 接口列表

### 1. 健康检查

```
GET /api/v1/health
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "ok",
    "db": "ok",
    "redis": "ok"
  }
}
```

### 2. 用户注册

```
POST /api/v1/auth/register
```

**请求体：**
```json
{
  "username": "demo",
  "email": "demo@example.com",
  "password": "123456"
}
```

**参数校验：**
- `username`: 必填，3-64 字符
- `email`: 必填，合法邮箱格式
- `password`: 必填，6-128 字符

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

**错误响应：**
```json
{
  "code": 40001,
  "message": "username already exists",
  "data": null
}
```

### 3. 用户登录

```
POST /api/v1/auth/login
```

**请求体：**
```json
{
  "username": "demo",
  "password": "123456"
}
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 7200
  }
}
```

### 4. 获取用户信息

```
GET /api/v1/auth/profile
Authorization: Bearer <access_token>
```

**成功响应：**
```json
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

### 5. 退出登录

```
POST /api/v1/auth/logout
Authorization: Bearer <access_token>
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```
