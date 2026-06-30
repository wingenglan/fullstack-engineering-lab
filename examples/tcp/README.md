# TCP 自定义协议演示

## 概述

本案例演示基于原始 TCP Socket 的自定义文本协议。Go 后端运行独立的 TCP Server，浏览器前端通过 HTTP API 管理 TCP 会话（创建/命令/关闭），命令响应通过 SSE 实时推送。

## 核心概念

- **TCP (Transmission Control Protocol)**：面向连接的、可靠的字节流传输协议，位于 OSI 第四层
- **自定义协议**：在 TCP 之上定义的应用层文本协议，命令行文本以换行符分隔
- **TCP vs HTTP**：HTTP 是建立在 TCP 之上的应用层协议（基于请求-响应模型），而原始 TCP 是持续的双向字节流，更灵活但需要自行处理粘包/拆包
- **会话管理**：每条 HTTP 创建的 TCP 会话对应一条独立的 TCP 连接，通过 SessionPool 管理生命周期
- **长连接**：TCP 连接在创建后保持活跃，可连续发送多条命令，直到主动关闭或超时

## 支持的命令

| 命令 | 参数 | 响应 | 说明 |
|------|------|------|------|
| `PING` | 无 | `PONG` | 心跳检测 |
| `ECHO` | `<msg>` | `<msg>` | 原样返回 |
| `UPPER` | `<msg>` | `<MSG>` | 转大写 |
| `LOWER` | `<msg>` | `<msg>` | 转小写 |
| `TIME` | 无 | RFC3339 时间戳 | 服务器当前时间 |
| `INFO` | 无 | 服务器信息 JSON | 版本、Go 版本等 |
| `HELP` | 无 | 命令列表 | 帮助文档 |
| `QUIT` | 无 | `BYE` | 断开 TCP 连接 |

## 快速开始

1. 启动项目：`make up`
2. 访问 TCP Demo：http://localhost:8888/cases/tcp
3. 创建会话 → 选择会话 → 输入命令（如 `PING`）→ 观察响应 → 关闭会话

## 技术栈

| 技术 | 用途 |
|------|------|
| Go + net | TCP Server 原生实现 |
| Gin | HTTP 框架 + SSE handler |
| Vue 3 | 前端框架 |
| EventSource (SSE) | 浏览器端 SSE 客户端 |
| Element Plus | UI 组件库 |

## 相关文档

- [API 文档](./api.md) — REST 接口 + SSE 推送格式说明
- [协议命令参考](./protocol.md) — 自定义协议完整命令手册
- [流程图](./flow.md) — 会话生命周期时序图
