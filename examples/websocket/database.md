# WebSocket 实时通讯 - 数据库设计

## 表结构

### chat_rooms - 聊天室表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键自增 |
| name | VARCHAR(128) | 房间名称 |
| description | VARCHAR(512) | 房间描述 |
| type | TINYINT | 类型：1=群聊，2=私聊 |
| creator_id | BIGINT | 创建者用户 ID |
| max_members | INT | 最大成员数（默认 500） |
| status | TINYINT | 状态：0=关闭，1=活跃 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |
| deleted_at | DATETIME | 软删除时间 |

**索引：**
- `idx_creator (creator_id)` - 按创建者查询
- `idx_deleted_at (deleted_at)` - 软删除过滤

### chat_messages - 聊天消息表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键自增 |
| room_id | BIGINT | 所属房间 ID |
| user_id | BIGINT | 发送者用户 ID |
| content | TEXT | 消息内容 |
| msg_type | TINYINT | 消息类型：1=文本，2=图片，3=系统 |
| created_at | DATETIME | 发送时间 |
| deleted_at | DATETIME | 软删除时间 |

**索引：**
- `idx_room_id (room_id)` - 按房间查询
- `idx_user_id (user_id)` - 按用户查询
- `idx_room_created (room_id, created_at)` - 房间消息时间排序（覆盖游标分页查询）
- `idx_deleted_at (deleted_at)` - 软删除过滤

## ER 关系

```
users (1) ──── (N) chat_messages
                    │
chat_rooms (1) ──── (N) chat_messages

users (1) ──── (N) chat_rooms (creator_id)
```

## 预置数据

```sql
INSERT INTO chat_rooms (name, description, type, creator_id) VALUES
  ('公共大厅', '所有用户的默认聊天室', 1, 0),
  ('技术交流', '技术讨论与问答', 1, 0),
  ('灌水区', '自由闲聊', 1, 0);
```

## 查询场景

### 获取房间历史消息（游标分页）

```sql
-- 获取 room_id=1 的最新 50 条消息
SELECT * FROM chat_messages
WHERE room_id = 1 AND deleted_at IS NULL
ORDER BY id DESC
LIMIT 51;

-- 获取 id=100 之前的 50 条消息
SELECT * FROM chat_messages
WHERE room_id = 1 AND id < 100 AND deleted_at IS NULL
ORDER BY id DESC
LIMIT 51;
```

> 多查 1 条用于判断 `has_more`。

### 注意事项

1. **游标分页 vs 偏移分页**：使用 `id < ?` 游标分页，避免 `OFFSET` 在大数据量下的性能问题
2. **复合索引**：`idx_room_created` 覆盖了最常见的查询模式（按房间查消息+时间排序）
3. **软删除**：消息不物理删除，通过 `deleted_at` 标记
4. **TEXT 类型**：消息内容使用 TEXT 类型，支持长文本
