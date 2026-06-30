# MQTT IoT 演示

## 概述

本案例演示基于 MQTT 协议的 IoT 传感器数据模拟，通过 Mosquitto Broker 发布传感器数据，后端订阅并桥接到 SSE（Server-Sent Events）推送至前端仪表盘实时展示。

## 核心概念

- **MQTT (Message Queuing Telemetry Transport)**：轻量级发布/订阅消息传输协议，专为低带宽、不可靠网络环境中的 IoT 设备设计
- **Broker（消息代理）**：MQTT 的核心组件，负责接收所有消息并路由给订阅者，本项目使用 Eclipse Mosquitto
- **Topic（主题）**：消息的分类标签，支持层级结构和通配符（`#` 匹配多级、`+` 匹配单级）
- **QoS（服务质量）**：0=最多一次（fire-and-forget）、1=至少一次（ACK 确认）、2=恰好一次（四步握手）
- **Pub/Sub 模型**：发布者与订阅者解耦，发布者无需知道谁在接收，订阅者也无需知道消息来源
- **SSE (Server-Sent Events)**：HTTP 长连接协议，服务端单向推送事件流到浏览器，比 WebSocket 更简单且支持自动重连

## 快速开始

1. 启动项目：`make up`
2. 访问 MQTT Demo：http://localhost:8888/cases/mqtt
3. 点击「启动模拟器」→ 观察仪表盘实时数据 → 手动发布消息 → 查看消息记录

## 技术栈

| 技术 | 用途 |
|------|------|
| Go + Gin | HTTP 框架 + SSE handler |
| Eclipse Paho (Go) | MQTT 客户端库 |
| Eclipse Mosquitto | MQTT Broker |
| Vue 3 | 前端框架 |
| EventSource (SSE) | 浏览器端 SSE 客户端 |
| Element Plus | UI 组件库 |

## 相关文档

- [API 文档](./api.md) — REST 接口 + SSE 推送格式说明
- [数据流时序图](./flow.md) — 设备→Broker→Server→SSE→前端的完整数据流
