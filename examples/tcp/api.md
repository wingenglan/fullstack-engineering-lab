# TCP 自定义协议演示 API 文档

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

---

## REST 接口列表

### 1. 创建 TCP 会话

```
POST /api/v1/tcp/sessions
```

无请求体。

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "session_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "created_at": "2026-06-30T10:00:00Z"
  }
}
```

---

### 2. 列出所有会话

```
GET /api/v1/tcp/sessions
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "session_id": "a1b2c3d4-...",
      "remote_addr": "127.0.0.1:9090",
      "created_at": "2026-06-30T10:00:00Z",
      "last_act_at": "2026-06-30T10:01:00Z",
      "duration_sec": 60,
      "command_count": 5,
      "is_alive": true
    }
  ]
}
```

---

### 3. 关闭会话

```
DELETE /api/v1/tcp/sessions/:id
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "session_id": "a1b2c3d4-...",
    "closed": true
  }
}
```

---

### 4. 发送命令

```
POST /api/v1/tcp/sessions/:id/send
```

**请求体：**
```json
{
  "command": "PING"
}
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "session_id": "a1b2c3d4-...",
    "command": "PING",
    "response": "PONG",
    "duration_ms": 2,
    "timestamp": "2026-06-30T10:00:05Z"
  }
}
```

**错误响应（命令无效或会话不存在）：**
```json
{
  "code": 50000,
  "message": "命令执行失败: ...",
  "data": null
}
```

---

### 5. 获取服务器统计

```
GET /api/v1/tcp/stats
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "server_addr": "0.0.0.0:9090",
    "active_sessions": 3,
    "max_sessions": 100,
    "uptime_sec": 3600
  }
}
```

---

### 6. 会话事件 SSE 流

```
GET /api/v1/tcp/sessions/:id/stream
```

**响应头：**
```
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
X-Accel-Buffering: no
```

**SSE 事件类型：**

#### `session_event` — 会话事件
```json
{
  "type": "connected",
  "session_id": "a1b2c3d4-...",
  "timestamp": "2026-06-30T10:00:00Z"
}
```

```json
{
  "type": "command_result",
  "session_id": "a1b2c3d4-...",
  "command": "PING",
  "response": "PONG",
  "duration_ms": 2,
  "timestamp": "2026-06-30T10:00:05Z"
}
```

```json
{
  "type": "disconnected",
  "session_id": "a1b2c3d4-...",
  "timestamp": "2026-06-30T10:05:00Z"
}
```

#### `heartbeat` — 保活心跳（每 30 秒）
```json
{
  "type": "heartbeat",
  "session_id": "a1b2c3d4-...",
  "timestamp": "2026-06-30T10:00:30Z"
}
```
