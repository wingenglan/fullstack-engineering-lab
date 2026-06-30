package mqtt

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

// SSEBroker SSE fan-out 事件分发器
type SSEBroker struct {
	clients    map[string]chan SSEEvent
	mu         sync.RWMutex
	history    []SSEEvent
	maxHistory int
}

// NewSSEBroker 创建 SSE 分发器
func NewSSEBroker(maxHistory int) *SSEBroker {
	return &SSEBroker{
		clients:    make(map[string]chan SSEEvent),
		history:    make([]SSEEvent, 0, maxHistory),
		maxHistory: maxHistory,
	}
}

// Subscribe 注册一个订阅者，返回只读通道和取消函数
func (b *SSEBroker) Subscribe(clientID string) (<-chan SSEEvent, func()) {
	ch := make(chan SSEEvent, 128)

	b.mu.Lock()
	b.clients[clientID] = ch
	b.mu.Unlock()

	unsub := func() {
		b.mu.Lock()
		delete(b.clients, clientID)
		b.mu.Unlock()
		close(ch)
	}

	return ch, unsub
}

// Publish 广播事件到所有订阅通道
func (b *SSEBroker) Publish(event SSEEvent) {
	// 追加到历史记录
	b.mu.Lock()
	if len(b.history) >= b.maxHistory {
		b.history = b.history[1:]
	}
	b.history = append(b.history, event)
	// 复制一份 clients 避免持有锁时写入通道
	clients := make(map[string]chan SSEEvent)
	for id, ch := range b.clients {
		clients[id] = ch
	}
	b.mu.Unlock()

	for clientID, ch := range clients {
		select {
		case ch <- event:
		default:
			zap.L().Debug("SSE 客户端通道满，丢弃事件", zap.String("client_id", clientID))
		}
	}
}

// History 返回历史消息快照
func (b *SSEBroker) History() []SSEEvent {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// 返回副本
	result := make([]SSEEvent, len(b.history))
	copy(result, b.history)
	return result
}

// ClientCount 订阅者数量
func (b *SSEBroker) ClientCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.clients)
}

// Cleanup 清理超时订阅者（心跳中调用）
func (b *SSEBroker) Cleanup(maxAge time.Duration) {
	// 简化实现：仅记录日志
	zap.L().Debug("SSE 当前订阅者数", zap.Int("count", b.ClientCount()))
}
