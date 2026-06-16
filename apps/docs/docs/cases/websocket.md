# WebSocket 实时通讯案例

## 概述

WebSocket 实时通讯案例演示了基于 WebSocket 协议的全双工实时通信，采用 Hub 模型管理连接，支持多房间聊天、消息持久化、在线状态和输入指示器，是一个完整的即时通讯系统实现。

## 工作原理

### WebSocket 协议

WebSocket 是一种在单个 TCP 连接上进行全双工通信的协议。与 HTTP 的请求-响应模式不同，WebSocket 允许服务端主动向客户端推送消息。

```
HTTP:   客户端 --请求--> 服务端 --响应--> 客户端
WebSocket: 客户端 <--双向--> 服务端
```

### 连接建立

WebSocket 连接通过 HTTP Upgrade 机制建立：

```
客户端: GET /api/v1/chat/ws?token=xxx
        Upgrade: websocket
        Connection: Upgrade

服务端: 101 Switching Protocols
        Upgrade: websocket
        Connection: Upgrade
```

### Hub 模型

服务端采用 Hub（集线器）模型管理所有 WebSocket 连接：

```
              ┌─────────┐
              │   Hub    │
              │ ┌──────┐ │
  Client A ──>│ │Room 1│ │<── Client B
              │ └──────┘ │
              │ ┌──────┐ │
  Client C ──>│ │Room 2│ │<── Client D
              │ └──────┘ │
              └─────────┘
```

- **Hub**：全局唯一的连接管理中心，负责注册/注销客户端、管理房间、广播消息
- **Room**：逻辑分组，同一房间内的客户端可以互相通信
- **Client**：每个 WebSocket 连接对应一个 Client，包含读写协程

### 核心特性

| 特性 | 说明 |
|------|------|
| 全双工通信 | 客户端和服务端可以随时互相发送消息 |
| 多房间支持 | 用户可以加入不同房间，消息只在房间内广播 |
| 消息持久化 | 消息同时写入 MySQL，新加入房间可查看历史消息 |
| 心跳保活 | 客户端定时 Ping，服务端 Pong，检测死连接 |
| 输入指示器 | 实时显示谁正在输入 |
| 自动重连 | 客户端断线后指数退避自动重连 |
| 在线状态 | 实时追踪每个房间的在线用户列表 |

## 消息协议

### 消息格式

所有 WebSocket 消息使用统一的 JSON 格式：

```json
{
  "type": "消息类型",
  "payload": { ... }
}
```

### 客户端 → 服务端

| 类型 | 说明 | Payload |
|------|------|---------|
| `join_room` | 加入房间 | `{ "room_id": 1 }` |
| `leave_room` | 离开房间 | `{ "room_id": 1 }` |
| `send_message` | 发送消息 | `{ "room_id": 1, "content": "你好", "msg_type": 1 }` |
| `typing` | 输入指示器 | `{ "room_id": 1 }` |
| `ping` | 心跳 | `{}` |

### 服务端 → 客户端

| 类型 | 说明 | Payload |
|------|------|---------|
| `new_message` | 新消息广播 | `{ "id", "room_id", "user_id", "username", "content", "msg_type", "created_at" }` |
| `user_joined` | 用户加入通知 | `{ "room_id", "user_id", "username" }` |
| `user_left` | 用户离开通知 | `{ "room_id", "user_id", "username" }` |
| `online_users` | 在线用户列表 | `{ "room_id", "users": [...], "count" }` |
| `user_typing` | 输入指示器 | `{ "room_id", "user_id", "username" }` |
| `room_history` | 历史消息 | `{ "room_id", "messages": [...] }` |
| `pong` | 心跳响应 | `{}` |
| `error` | 错误消息 | `{ "code", "message" }` |

## REST 接口

### GET /api/v1/chat/rooms

获取所有可用的聊天室及其在线人数。

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "公共大厅",
      "description": "所有用户的默认聊天室",
      "type": 1,
      "creator_id": 0,
      "member_count": 3,
      "status": 1,
      "created_at": "2026-06-16 10:00:00"
    }
  ]
}
```

### POST /api/v1/chat/rooms

创建聊天室（需认证）。

```json
// 请求
{
  "name": "前端技术交流",
  "description": "讨论前端技术",
  "type": 1
}

// 响应
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 4,
    "name": "前端技术交流",
    "description": "讨论前端技术",
    "type": 1,
    "creator_id": 1,
    "member_count": 0,
    "status": 1,
    "created_at": "2026-06-16 12:00:00"
  }
}
```

### GET /api/v1/chat/messages?room_id=1&limit=50

获取消息历史（游标分页）。

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "messages": [
      {
        "id": 1,
        "room_id": 1,
        "user_id": 1,
        "username": "admin",
        "content": "欢迎来到聊天室！",
        "msg_type": 1,
        "created_at": "2026-06-16 10:00:00"
      }
    ],
    "has_more": false
  }
}
```

## 代码架构

```
apps/server/
├── pkg/websocket/
│   ├── hub.go          # Hub 连接管理中心
│   ├── client.go       # 客户端连接（读写协程）
│   ├── message.go      # 消息协议定义
│   └── upgrader.go     # HTTP 升级器
├── internal/
│   ├── handler/
│   │   ├── chat_handler.go  # REST API 处理器
│   │   └── ws_handler.go    # WebSocket 连接处理
│   ├── service/
│   │   └── chat_service.go  # 业务逻辑
│   ├── repository/
│   │   └── chat_repo.go     # 数据层
│   └── model/
│       ├── chat.go          # 数据模型
│       └── dto.go           # 请求/响应 DTO
```

```
apps/web/
├── src/
│   ├── composables/
│   │   └── useWebSocket.ts  # WebSocket 连接 composable
│   ├── stores/
│   │   └── chat.ts          # Pinia 聊天状态管理
│   ├── api/
│   │   └── chat.ts          # REST API 封装
│   ├── views/cases/
│   │   └── WebSocketView.vue # 聊天室页面
│   └── types/
│       └── index.ts         # TypeScript 类型定义
```

## 适用场景

- **在线客服系统**：客户与客服实时沟通
- **协同编辑**：多人同时编辑文档
- **实时通知**：订单状态变更、系统告警
- **在线教育**：课堂互动、弹幕
- **多人游戏**：游戏状态同步

## 注意事项

1. **认证方式**：WebSocket 通过 URL 参数传递 Token（浏览器不支持自定义 Header）
2. **消息大小**：限制单条消息最大 64KB，防止滥用
3. **连接数**：每个用户支持多标签页（多连接），通过 userID 关联
4. **内存管理**：空房间自动清理，防止内存泄漏
5. **生产环境**：应使用 Redis Pub/Sub 实现多实例消息广播
