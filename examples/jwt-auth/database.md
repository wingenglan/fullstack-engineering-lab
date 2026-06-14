# JWT Auth 数据库设计

## Users 表

### 表结构

```sql
CREATE TABLE users (
  id            BIGINT PRIMARY KEY AUTO_INCREMENT,
  username      VARCHAR(64)  NOT NULL UNIQUE,
  email         VARCHAR(128) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  nickname      VARCHAR(64)  DEFAULT '',
  status        TINYINT      NOT NULL DEFAULT 1,
  created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at    DATETIME     NULL
);
```

### 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键，自增 |
| username | VARCHAR(64) | 用户名，唯一 |
| email | VARCHAR(128) | 邮箱，唯一 |
| password_hash | VARCHAR(255) | bcrypt 哈希后的密码 |
| nickname | VARCHAR(64) | 昵称，可选 |
| status | TINYINT | 状态：0=禁用，1=启用 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间（自动） |
| deleted_at | DATETIME | 软删除时间（GORM） |

### 索引

- `idx_username`: username 字段索引
- `idx_email`: email 字段索引
- `idx_deleted_at`: 软删除字段索引

### GORM 模型映射

```go
type User struct {
    gorm.Model
    Username     string `gorm:"type:varchar(64);uniqueIndex;not null"`
    Email        string `gorm:"type:varchar(128);uniqueIndex;not null"`
    PasswordHash string `gorm:"type:varchar(255);not null"`
    Nickname     string `gorm:"type:varchar(64);default:''"`
    Status       int8   `gorm:"type:tinyint;not null;default:1"`
}
```

### 密码存储说明

密码使用 bcrypt 算法进行哈希处理：
- 算法: bcrypt
- Cost: 10
- 输出长度: 60 字符（固定）
- 存储格式: `$2a$10$...`

**绝不存储明文密码。**
