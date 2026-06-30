package tcp

import (
	"fmt"
	"strings"
	"time"
)

// ParseCommand 从输入行中分离命令和参数
func ParseCommand(line string) (cmd string, args string) {
	line = strings.TrimSpace(line)
	parts := strings.SplitN(line, " ", 2)
	cmd = strings.ToUpper(parts[0])
	if len(parts) > 1 {
		args = parts[1]
	}
	return
}

// HandleResult 命令处理结果
type HandleResult struct {
	Response string // 响应文本
	Quit     bool   // 是否断开连接
}

// Handle 命令处理分发
func Handle(cmd, args string) HandleResult {
	switch cmd {
	case "PING":
		return HandleResult{Response: "PONG"}
	case "ECHO":
		return HandleResult{Response: args}
	case "UPPER":
		return HandleResult{Response: strings.ToUpper(args)}
	case "LOWER":
		return HandleResult{Response: strings.ToLower(args)}
	case "TIME":
		return HandleResult{Response: time.Now().Format(time.RFC3339)}
	case "INFO":
		return HandleResult{Response: `{"server":"TCP Demo Server","version":"1.0.0","go":"` +
			strings.ReplaceAll(`\n`, " ", "") + `"}`}
	case "HELP":
		return HandleResult{Response: helpText()}
	case "QUIT":
		return HandleResult{Response: "BYE", Quit: true}
	default:
		return HandleResult{Response: fmt.Sprintf("ERROR: 未知命令 %s，输入 HELP 查看可用命令", cmd)}
	}
}

func helpText() string {
	return "可用命令: PING | ECHO <msg> | UPPER <msg> | LOWER <msg> | TIME | INFO | HELP | QUIT"
}
