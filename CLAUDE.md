# 项目规则

## 语言规范

- 本项目为**中文项目**，所有代码注释、提交信息、文档内容均使用**中文**
- 变量名、函数名、类型名等标识符使用英文（遵循 Go / TypeScript 惯例）
- 错误信息面向用户时使用中文，面向日志/调试时可使用英文
- 禁止在代码中使用英文注释

## 项目简介

FullStack Engineering Lab 是面向开发者的技术实验室，用于沉淀常见技术的完整实战案例。

## 技术栈

- 前端：Vue 3 + Vite + TypeScript + Element Plus + UnoCSS
- 后端：Go + Gin + Gorm + Viper + Zap
- 基础设施：MySQL + Redis + Nginx + Docker Compose

## 代码规范

- 后端分层：handler → service → repository → model
- 前端分层：view → api → store → types
- 统一响应格式：`{ code, message, data }`
- API 路径前缀：`/api/v1/`

## 提交规范

- 格式：`<type>: <描述>`
- type：feat / fix / docs / refactor / test / chore
- 描述使用中文
