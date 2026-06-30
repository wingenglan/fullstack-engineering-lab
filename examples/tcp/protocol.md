# TCP 自定义协议 - 命令参考手册

## 协议规范

- **传输层**：TCP（面向连接、可靠字节流）
- **消息格式**：文本行，UTF-8 编码，以 `\n` (换行) 结尾
- **命令大小写**：不敏感，服务端自动转换为大写处理
- **响应格式**：单行文本，以 `\n` 结尾
- **连接模型**：长连接，可连续发送多条命令

---

## 命令列表

### PING — 心跳检测

检测连接是否存活。

```
> PING
< PONG
```

### ECHO — 原样返回

将参数原样返回。

```
> ECHO hello world
< hello world
```

### UPPER — 转大写

将参数转换为大写后返回。

```
> UPPER hello world
< HELLO WORLD
```

### LOWER — 转小写

将参数转换为小写后返回。

```
> LOWER HELLO WORLD
< hello world
```

### TIME — 获取时间

返回服务器当前时间，RFC3339 格式。

```
> TIME
< 2026-06-30T10:00:00+08:00
```

### INFO — 获取服务器信息

返回服务器信息 JSON。

```
> INFO
< {"server":"TCP Demo Server","version":"1.0.0"}
```

### HELP — 获取帮助

返回所有可用命令列表。

```
> HELP
< 可用命令: PING | ECHO <msg> | UPPER <msg> | LOWER <msg> | TIME | INFO | HELP | QUIT
```

### QUIT — 断开连接

关闭 TCP 连接。

```
> QUIT
< BYE
（连接关闭）
```

---

## 错误处理

### 未知命令

```
> UNKNOWN_CMD
< ERROR: 未知命令 UNKNOWN_CMD，输入 HELP 查看可用命令
```

### 连接超时

- TCP Server 默认读超时 30 秒，写超时 10 秒
- 超时后连接自动关闭
- 前端 SSE 流会收到 `disconnected` 事件

---

## 粘包/拆包处理

本协议使用**换行分隔**方式解决 TCP 粘包/拆包问题：

- 服务端一次 Read 可能读到多个命令（粘包）
- 服务端一次 Read 可能只读到不完整的命令（拆包）

当前简化实现中，每次 `SendCommand` 对应一次 Write + 一次 Read，避免粘包问题。生产环境建议使用 `bufio.Scanner` 逐行读取。
