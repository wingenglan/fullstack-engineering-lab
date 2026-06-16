# WebSocket 实时通讯 API 文档

## 基础信息

- Base URL: `/api/v1`
- WebSocket URL: `ws://host/api/v1/chat/ws?token=<jwt_token>`
- Content-Type: `application/json`
- 认证方式:
  - REST API: Bearer Token（部分接口公开）
  - WebSocket: URL 参数 `?token=<jwt_token>`

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
| 40010 | 聊天错误（参数无效等） |
| 40011 | 房间已满 |
| 50000 | 服务器内部错误 |

---

## REST 接口列表

### 1. 获取聊天室列表

```
GET /api/v1/chat/rooms
```

**认证：** 无需

**响应：**
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
    },
    {
      "id": 2,
      "name": "技术交流",
      "description": "技术讨论与问答",
      "type": 1,
      "creator_id": 0,
      "member_count": 1,
      "status": 1,
      "created_at": "2026-06-16 10:00:00"
    }
  ]
}
```

---

### 2. 创建聊天室

```
POST /api/v1/chat/rooms
```

**认证：** Bearer Token

**请求体：**
```json
{
  "name": "前端技术交流",
  "description": "讨论前端技术",
  "type": 1
}
```

**参数说明：**
- `name`: 必填，房间名称（2-128 字符）
- `description`: 可选，房间描述（最大 512 字符）
- `type`: 必填，房间类型（1=群聊，2=私聊）

**成功响应：**
```json
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

---

### 3. 获取聊天室详情

```
GET /api/v1/chat/rooms/:id
```

**认证：** 无需

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "公共大厅",
    "description": "所有用户的默认聊天室",
    "type": 1,
    "creator_id": 0,
    "member_count": 3,
    "status": 1,
    "created_at": "2026-06-16 10:00:00"
  }
}
```

**房间不存在：**
```json
{
  "code": 40010,
  "message": "聊天室不存在",
  "data": null
}
```

---

### 4. 获取消息历史

```
GET /api/v1/chat/messages?room_id=1&limit=50&before=100
```

**认证：** 无需

**查询参数：**
- `room_id`: 必填，聊天室 ID
- `limit`: 可选，每页数量（默认 50，最大 100）
- `before`: 可选，获取此 ID 之前的消息（游标分页）

**成功响应：**
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

---

### 5. 获取在线用户

```
GET /api/v1/chat/rooms/:id/online
```

**认证：** 无需

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    { "user_id": 1, "username": "admin" },
    { "user_id": 2, "username": "user1" }
  ]
}
```

---

### 6. WebSocket 连接

```
GET /api/v1/chat/ws?token=<jwt_token>
Upgrade: websocket
Connection: Upgrade
```

**认证：** URL 参数 Token

**说明：**
- 客户端通过 HTTP GET 请求升级为 WebSocket 连接
- Token 通过 URL 参数传递（浏览器 WebSocket API 不支持自定义 Header）
- 连接成功后进入消息模式，使用 JSON 格式通信

---

## WebSocket 消息协议

### 消息格式

```json
{
  "type": "消息类型",
  "payload": { ... }
}
```

### 客户端 → 服务端

#### join_room - 加入房间

```json
{
  "type": "join_room",
  "payload": { "room_id": 1 }
}
```

服务端响应：
1. `room_history` - 推送最近 50 条历史消息
2. `user_joined` - 广播加入通知（发给房间其他人）
3. `online_users` - 推送更新后的在线用户列表

#### leave_room - 离开房间

```json
{
  "type": "leave_room",
  "payload": { "room_id": 1 }
}
```

服务端响应：
1. `user_left` - 广播离开通知
2. `online_users` - 推送更新后的在线用户列表

#### send_message - 发送消息

```json
{
  "type": "send_message",
  "payload": {
    "room_id": 1,
    "content": "大家好！",
    "msg_type": 1
  }
}
```

**payload 参数：**
- `room_id`: 必填，目标房间 ID
- `content`: 必填，消息内容（最大 5000 字符）
- `msg_type`: 可选，消息类型（1=文本，2=图片，3=系统，默认 1）

服务端响应：`new_message` 广播到房间所有人（包括发送者）

#### typing - 输入指示器

```json
{
  "type": "typing",
  "payload": { "room_id": 1 }
}
```

服务端响应：`user_typing` 广播给房间其他人（排除发送者）

#### ping - 心跳

```json
{
  "type": "ping",
  "payload": {}
}
```

服务端响应：`pong`

### 服务端 → 客户端

#### new_message - 新消息

```json
{
  "type": "new_message",
  "payload": {
    "id": 42,
    "room_id": 1,
    "user_id": 1,
    "username": "admin",
    "content": "大家好！",
    "msg_type": 1,
    "created_at": "2026-06-16T12:00:00+08:00"
  }
}
```

#### user_joined - 用户加入

```json
{
  "type": "user_joined",
  "payload": {
    "room_id": 1,
    "user_id": 2,
    "username": "user1"
  }
}
```

#### user_left - 用户离开

```json
{
  "type": "user_left",
  "payload": {
    "room_id": 1,
    "user_id": 2,
    "username": "user1"
  }
}
```

#### online_users - 在线用户列表

```json
{
  "type": "online_users",
  "payload": {
    "room_id": 1,
    "users": [
      { "user_id": 1, "username": "admin" },
      { "user_id": 2, "username": "user1" }
    ],
    "count": 2
  }
}
```

#### user_typing - 输入指示器

```json
{
  "type": "user_typing",
  "payload": {
    "room_id": 1,
    "user_id": 2,
    "username": "user1"
  }
}
```

#### room_history - 历史消息

```json
{
  "type": "room_history",
  "payload": {
    "room_id": 1,
    "messages": [
      {
        "id": 1,
        "user_id": 1,
        "username": "admin",
        "content": "欢迎！",
        "msg_type": 1,
        "created_at": "2026-06-16 10:00:00"
      }
    ]
  }
}
```

#### error - 错误

```json
{
  "type": "error",
  "payload": {
    "code": 400,
    "message": "消息内容不能为空"
  }
}
```
