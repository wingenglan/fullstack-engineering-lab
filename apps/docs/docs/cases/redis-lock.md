# Redis 分布式锁案例

## 概述

Redis 分布式锁案例演示了基于 Redis 的分布式锁实现，使用 `SET NX EX` 命令实现加锁、Lua 脚本实现原子释放，适用于分布式系统中的并发控制场景。

## 工作原理

### 加锁（SET NX EX）

```redis
SET lock:resource <unique_token> NX EX 10
```

- **NX**：仅当 key 不存在时才设置（互斥保证）
- **EX 10**：设置 10 秒过期时间（防死锁）
- **unique_token**：随机生成的 UUID，标识锁持有者

### 释放锁（Lua 脚本）

```lua
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end
```

使用 Lua 脚本保证「比较 value + 删除」是原子操作，防止误删其他客户端持有的锁。

### 核心特性

| 特性 | 说明 |
|------|------|
| 互斥性 | NX 保证同一时刻只有一个客户端能获取锁 |
| 安全性 | Owner Token + Lua 脚本保证只有持有者能释放锁 |
| 活锁防护 | TTL 自动过期，即使持有者崩溃锁也会释放 |
| 原子性 | 加锁和释放操作都是原子的 |

## 接口文档

### POST /api/v1/lock/acquire

获取分布式锁。

```json
// 请求
{
  "resource": "order:10086",
  "ttl": 10
}

// 成功响应
{
  "code": 0,
  "message": "success",
  "data": {
    "resource": "order:10086",
    "owner": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "ttl": 10
  }
}

// 锁冲突响应
{
  "code": 40009,
  "message": "resource \"order:10086\" is locked by a1b2c3d4 (TTL: 8500ms)",
  "data": null
}
```

### POST /api/v1/lock/release

释放分布式锁。

```json
// 请求
{
  "resource": "order:10086"
}

// 响应
{
  "code": 0,
  "message": "lock released",
  "data": null
}
```

### POST /api/v1/lock/status

查询锁状态。

```json
// 请求
{
  "resource": "order:10086"
}

// 响应
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

### POST /api/v1/lock/contention

并发争抢演示——模拟多个协程同时争抢同一把锁。

```json
// 请求
{
  "resource": "demo-resource",
  "ttl": 10,
  "goroutines": 5,
  "hold_ms": 500
}

// 响应
{
  "code": 0,
  "message": "success",
  "data": {
    "resource": "demo-resource",
    "results": [
      { "goroutine_id": 1, "acquired": true, "wait_ms": 0, "message": "获取锁成功，持有 500ms 后释放" },
      { "goroutine_id": 2, "acquired": false, "wait_ms": 0, "message": "获取锁失败：资源已被占用" },
      { "goroutine_id": 3, "acquired": false, "wait_ms": 0, "message": "获取锁失败：资源已被占用" }
    ],
    "summary": { "total": 5, "succeeded": 1, "failed": 4 }
  }
}
```

## 代码架构

```
apps/server/
├── pkg/redislock/
│   └── lock.go              # 核心分布式锁实现
├── internal/
│   ├── handler/
│   │   └── lock_handler.go  # HTTP 处理器
│   ├── service/
│   │   └── lock_service.go  # 业务逻辑
│   └── model/
│       └── dto.go           # 请求/响应 DTO
```

## 适用场景

- **秒杀扣库存**：防止超卖
- **订单防重**：防止重复下单
- **定时任务互斥**：多实例部署时同一任务只执行一次
- **资源配额控制**：限制并发访问数

## 注意事项

1. **锁粒度**：资源名应尽量细粒度，避免全局锁
2. **TTL 选择**：应大于业务执行时间，但不宜过长
3. **时钟漂移**：Redis 主从切换可能导致锁丢失（Redlock 可解决）
4. **可重入**：本实现为不可重入锁，同一客户端重复获取会失败
