# Redis Lock API 文档

## 基础信息

- Base URL: `/api/v1`
- Content-Type: `application/json`
- 认证方式: 无（演示用途，所有接口公开访问）

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
| 40009 | 锁冲突（资源已被占用） |
| 50000 | 服务器内部错误 |

---

## 接口列表

### 1. 获取分布式锁

```
POST /api/v1/lock/acquire
```

**请求体：**
```json
{
  "resource": "order:10086",
  "ttl": 10
}
```

**参数说明：**
- `resource`: 必填，资源名称（将作为 Redis key 的一部分）
- `ttl`: 必填，锁过期时间（秒），范围 1-300

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "resource": "order:10086",
    "owner": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "ttl": 10
  }
}
```

**锁冲突响应：**
```json
{
  "code": 40009,
  "message": "resource \"order:10086\" is locked by a1b2c3d4 (TTL: 8500ms)",
  "data": null
}
```

---

### 2. 释放分布式锁

```
POST /api/v1/lock/release
```

**请求体：**
```json
{
  "resource": "order:10086"
}
```

**成功响应：**
```json
{
  "code": 0,
  "message": "lock released",
  "data": null
}
```

---

### 3. 查询锁状态

```
POST /api/v1/lock/status
```

**请求体：**
```json
{
  "resource": "order:10086"
}
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "resource": "order:10086",
    "locked": true,
    "owner": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "ttl_ms": 8500
  }
}
```

---

### 4. 并发争抢演示

```
POST /api/v1/lock/contention
```

模拟多个 Go 协程同时争抢同一把分布式锁。

**请求体：**
```json
{
  "resource": "demo-resource",
  "ttl": 10,
  "goroutines": 5,
  "hold_ms": 500
}
```

**参数说明：**
- `resource`: 必填，争抢的资源名
- `ttl`: 必填，锁 TTL（秒）
- `goroutines`: 必填，并发协程数（2-20）
- `hold_ms`: 必填，持有锁的模拟时间（100-10000ms）

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "resource": "demo-resource",
    "results": [
      { "goroutine_id": 1, "acquired": true, "wait_ms": 0, "message": "获取锁成功，持有 500ms 后释放" },
      { "goroutine_id": 2, "acquired": false, "wait_ms": 0, "message": "获取锁失败：资源已被占用" },
      { "goroutine_id": 3, "acquired": false, "wait_ms": 0, "message": "获取锁失败：资源已被占用" },
      { "goroutine_id": 4, "acquired": false, "wait_ms": 0, "message": "获取锁失败：资源已被占用" },
      { "goroutine_id": 5, "acquired": false, "wait_ms": 0, "message": "获取锁失败：资源已被占用" }
    ],
    "summary": {
      "total": 5,
      "succeeded": 1,
      "failed": 4
    }
  }
}
```
