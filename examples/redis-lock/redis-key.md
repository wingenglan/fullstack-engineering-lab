# Redis 分布式锁 - Redis Key 设计

## Key 命名规范

```
lock:{resource}
```

### 示例

| 资源 | Redis Key | 说明 |
|------|-----------|------|
| order:10086 | `lock:order:10086` | 订单锁 |
| user:login | `lock:user:login` | 登录锁 |
| payment:9527 | `lock:payment:9527` | 支付锁 |
| cron:daily-report | `lock:cron:daily-report` | 定时任务锁 |

## Value 设计

Value 是一个随机生成的 UUID 格式字符串（如 `a1b2c3d4-e5f6-7890-abcd-ef1234567890`），用于标识锁的持有者。

**为什么需要 Owner Token？**

```
场景：客户端 A 获取了锁，但执行时间超过 TTL，锁自动过期
      客户端 B 此时获取了同一把锁
      客户端 A 执行完毕，尝试释放锁

如果没有 Owner 校验：A 会误删 B 的锁
有 Owner 校验：A 发现 value 不匹配，释放失败
```

## TTL 策略

| 场景 | 建议 TTL | 说明 |
|------|---------|------|
| 短时操作（API 请求） | 5-10s | 请求超时一般 < 5s |
| 业务处理（下单等） | 10-30s | 需要数据库操作 |
| 定时任务 | 60-300s | 任务执行时间较长 |
| 批量操作 | 30-120s | 视数据量而定 |

## Redis 命令交互

### 加锁

```redis
> SET lock:order:10086 a1b2c3d4-e5f6-7890 NX EX 10
OK                          # 获取成功
(nil)                       # 获取失败（已存在）
```

### 查看锁状态

```redis
> GET lock:order:10086
"a1b2c3d4-e5f6-7890-abcd-ef1234567890"

> TTL lock:order:10086
(integer) 7                 # 剩余 7 秒
```

### 释放锁（Lua 脚本）

```redis
> EVAL "if redis.call('GET',KEYS[1])==ARGV[1] then return redis.call('DEL',KEYS[1]) else return 0 end" 1 lock:order:10086 a1b2c3d4-e5f6-7890
(integer) 1                 # 释放成功
(integer) 0                 # 释放失败（value 不匹配）
```

## 注意事项

1. **Key 不要使用全局锁**：如 `lock:global`，会成为性能瓶颈
2. **合理设置 TTL**：TTL 应大于业务最大执行时间
3. **避免 Key 膨胀**：锁释放后 key 会自动删除（TTL 过期或主动 DEL）
4. **Redis 内存策略**：建议使用 `allkeys-lru`（已在 redis.conf 中配置）
