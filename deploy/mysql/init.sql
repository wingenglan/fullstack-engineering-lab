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
