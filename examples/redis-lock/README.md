# Redis 分布式锁 Demo

## 概述

本案例演示基于 Redis 的分布式锁完整实现，包括加锁、释放锁、锁状态查询、并发争抢模拟等功能。通过真实可运行的接口，直观展示分布式环境下的互斥控制。

## 核心概念

- **分布式锁**：多进程/多节点环境下对共享资源的互斥访问控制
- **SET NX EX**：Redis 原子命令，`NX`（Not Exist）保证只有一个客户端能成功写入
- **Owner Token**：随机 UUID，标识锁的持有者，防止误删他人的锁
- **Lua 脚本**：用于释放锁的原子操作，确保"比对 + 删除"不被打断
- **TTL 防死锁**：锁设置过期时间，即使持有者崩溃也能自动释放

## 快速开始

1. 启动项目：`make up`
2. 访问 Redis Lock Demo：http://localhost:3000/cases/redis-lock
3. 获取锁 → 查询状态 → 释放锁 → 模拟并发争抢

## 技术栈

| 技术 | 用途 |
|------|------|
| Go + Gin | HTTP 框架 |
| go-redis/v9 | Redis 客户端 |
| Redis SET NX EX | 原子加锁命令 |
| Redis EVAL (Lua) | 原子解锁脚本 |
| Vue 3 | 前端框架 |
| Element Plus | UI 组件库 |

## 相关文档

- [API 文档](./api.md) — 接口说明与请求/响应示例
- [Redis Key 设计](./redis-key.md) — Key 命名规范、Value 设计、TTL 策略
- [流程图](./flow.md) — 加锁/解锁时序图与并发争抢流程
