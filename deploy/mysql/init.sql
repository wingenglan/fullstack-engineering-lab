-- FullStack Engineering Lab - Database Init
-- This script runs automatically when MySQL container starts for the first time.

CREATE DATABASE IF NOT EXISTS engineering_lab
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE engineering_lab;

-- Users table for JWT Auth Demo
CREATE TABLE IF NOT EXISTS users (
  id           BIGINT PRIMARY KEY AUTO_INCREMENT,
  username     VARCHAR(64)  NOT NULL UNIQUE,
  email        VARCHAR(128) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  nickname     VARCHAR(64)  DEFAULT '',
  status       TINYINT      NOT NULL DEFAULT 1 COMMENT '0=disabled, 1=active',
  created_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at   DATETIME     NULL,
  INDEX idx_username (username),
  INDEX idx_email (email),
  INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Chat rooms table for WebSocket Real-time Chat
CREATE TABLE IF NOT EXISTS chat_rooms (
  id          BIGINT PRIMARY KEY AUTO_INCREMENT,
  name        VARCHAR(128) NOT NULL,
  description VARCHAR(512) DEFAULT '',
  type        TINYINT      NOT NULL DEFAULT 1 COMMENT '1=group, 2=private',
  creator_id  BIGINT       NOT NULL,
  max_members INT          NOT NULL DEFAULT 500,
  status      TINYINT      NOT NULL DEFAULT 1 COMMENT '0=closed, 1=active',
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at  DATETIME     NULL,
  INDEX idx_creator (creator_id),
  INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Chat messages table for WebSocket Real-time Chat
CREATE TABLE IF NOT EXISTS chat_messages (
  id          BIGINT PRIMARY KEY AUTO_INCREMENT,
  room_id     BIGINT       NOT NULL,
  user_id     BIGINT       NOT NULL,
  content     TEXT         NOT NULL,
  msg_type    TINYINT      NOT NULL DEFAULT 1 COMMENT '1=text, 2=image, 3=system',
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at  DATETIME     NULL,
  INDEX idx_room_id (room_id),
  INDEX idx_user_id (user_id),
  INDEX idx_room_created (room_id, created_at),
  INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert default chat rooms
INSERT INTO chat_rooms (name, description, type, creator_id) VALUES
  ('公共大厅', '所有用户的默认聊天室', 1, 0),
  ('技术交流', '技术讨论与问答', 1, 0),
  ('灌水区', '自由闲聊', 1, 0);
