# JWT 认证授权 Demo

## 概述

本案例演示基于 JSON Web Token (JWT) 的完整认证授权流程，包括用户注册、登录、Token 签发与验证、Token 黑名单等功能。

## 核心概念

- **JWT (JSON Web Token)**: 一种开放标准 (RFC 7519)，用于在各方之间安全地传输信息
- **Access Token**: 访问令牌，用于访问受保护的资源
- **Bearer Token**: HTTP 请求头中的认证方式
- **bcrypt**: 密码哈希算法，保证密码安全存储
- **Token Blacklist**: 使用 Redis 存储已注销的 Token

## 快速开始

1. 启动项目：`make up`
2. 访问 JWT Demo：http://localhost:3000/cases/jwt-auth
3. 注册新用户 → 登录 → 查看 Token → 获取用户信息 → 退出登录

## 技术栈

| 技术 | 用途 |
|------|------|
| Go + Gin | HTTP 框架 |
| Gorm | ORM，操作 MySQL |
| golang-jwt/jwt/v5 | JWT 签发与验证 |
| golang.org/x/crypto/bcrypt | 密码哈希 |
| go-redis/v9 | Redis 客户端，Token 黑名单 |
| Vue 3 | 前端框架 |
| Element Plus | UI 组件库 |
| Pinia | 状态管理 |
