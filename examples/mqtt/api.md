# MQTT IoT 演示 API 文档

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

### 1. 推出消息

```
POST /api/v1/mqtt/publish
```

**请求体：**
```json
{
  "topic": "sensors/test",
  "payload": "{\"value\": 42}",
  "qos": 1
}
```

**参数说明：**
- `topic`: 必填，MQTT Topic 名称
- `payload`: 必填，消息载荷（合法 JSON 字符串）
- `qos`: 可选，服务质量等级（0/1/2，默认 1）

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "topic": "sensors/test",
    "timestamp": "2026-06-30T10:00:00Z"
  }
}
```

---

### 2. 获取消息记录

```
GET /api/v1/mqtt/messages?limit=50
```

**查询参数：**
- `limit`: 可选，返回最近 N 条消息（默认 50，最大 100）

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "topic": "sensors/temperature",
      "payload": "{\"device_id\":\"sensor-01\",\"value\":23.5,\"unit\":\"°C\",\"sensor_type\":\"temperature\",\"timestamp\":\"...\"}",
      "source": "simulator",
      "timestamp": "2026-06-30T10:00:00Z"
    }
  ]
}
```

---

### 3. 启动模拟器

```
POST /api/v1/mqtt/simulator/start
```

**请求体（可选）：**
```json
{
  "interval_ms": 2000
}
```

**参数说明：**
- `interval_ms`: 可选，发布间隔毫秒（默认使用配置文件值 2000）

---

### 4. 停止模拟器

```
POST /api/v1/mqtt/simulator/stop
```

无请求体。

---

### 5. 查询模拟器状态

```
GET /api/v1/mqtt/simulator/status
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "running": true,
    "interval_ms": 2000,
    "message_count": 42
  }
}
```

---

### 6. SSE 订阅

```
GET /api/v1/mqtt/subscribe
```

**响应头：**
```
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
X-Accel-Buffering: no
```

**SSE 事件类型：**

#### `connected` — 连接确认
```json
{
  "client_id": "uuid-xxx",
  "history_count": 10
}
```

#### `mqtt_message` — MQTT 消息到达
```json
{
  "topic": "sensors/temperature",
  "payload": "{\"device_id\":\"sensor-01\",\"value\":23.5,\"unit\":\"°C\",\"sensor_type\":\"temperature\",\"timestamp\":\"...\"}",
  "source": "simulator",
  "timestamp": "2026-06-30T10:00:00Z"
}
```

#### `heartbeat` — 保活心跳（每 30 秒）
```json
{"timestamp": "2026-06-30T10:00:30Z"}
```
