# WebSocket 实时通讯 Demo

## 概述

本案例演示基于 WebSocket 的实时多人聊天室，包括聊天室管理、实时消息收发、在线用户同步、输入指示器、心跳保活与断线重连等功能。

## 核心概念

- **WebSocket**：基于 TCP 的全双工通信协议，一次握手后服务端可主动推送消息
- **Hub（连接中枢）**：服务端维护的连接管理器，负责消息广播与客户端注册/注销
- **消息协议**：自定义 JSON 格式 `{ type, payload }`，通过消息类型区分不同事件
- **心跳保活**：客户端每 30 秒发送 `ping`，服务端回复 `pong`，防止连接被中间设备超时断开
- **指数退避重连**：断线后按 `min(1000 * 2^n, 30000) + jitter` 策略自动重连，最多重试 10 次
- **游标分页**：历史消息使用 `id < ?` 游标分页，避免大数据量下 OFFSET 性能问题

## 快速开始

1. 启动项目：`make up`
2. 访问 WebSocket Demo：http://localhost:3000/cases/websocket
3. 登录账号 → 建立 WebSocket 连接 → 加入聊天室 → 发送消息

> WebSocket 连接需要先登录获取 JWT Token，Token 通过 URL 参数传递：`ws://host/api/v1/chat/ws?token=<jwt_token>`

## 技术栈

| 技术 | 用途 |
|------|------|
| Go + Gin | HTTP / WebSocket 升级 |
| gorilla/websocket | WebSocket 服务端实现 |
| Gorm | ORM，持久化聊天消息 |
| MySQL | 存储聊天室与消息记录 |
| Vue 3 | 前端框架 |
| Element Plus | UI 组件库 |
| Pinia | 状态管理（在线用户、消息列表） |

## 相关文档

- [API 文档](./api.md) — REST 接口 + WebSocket 消息协议说明
- [数据库设计](./database.md) — 聊天室表、消息表结构与索引设计
- [连接流程](./flow.md) — 完整连接时序图与断线重连策略
