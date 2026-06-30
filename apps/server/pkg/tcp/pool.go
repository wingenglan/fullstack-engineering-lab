package tcp

import (
	"fmt"
	"sync"
)

// SessionPool 并发安全的会话池
type SessionPool struct {
	sessions map[string]*Session
	mu       sync.RWMutex
	maxSize  int
}

// NewSessionPool 创建会话池
func NewSessionPool(maxSize int) *SessionPool {
	return &SessionPool{
		sessions: make(map[string]*Session),
		maxSize:  maxSize,
	}
}

// Add 添加会话
func (p *SessionPool) Add(s *Session) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.sessions) >= p.maxSize {
		return fmt.Errorf("会话池已满（最大 %d）", p.maxSize)
	}
	p.sessions[s.ID] = s
	return nil
}

// Get 获取会话
func (p *SessionPool) Get(id string) (*Session, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	s, ok := p.sessions[id]
	return s, ok
}

// Remove 从池中移除（不关闭连接）
func (p *SessionPool) Remove(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.sessions, id)
}

// List 列出所有会话
func (p *SessionPool) List() []*SessionStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	stats := make([]*SessionStats, 0, len(p.sessions))
	for _, s := range p.sessions {
		stats = append(stats, s.Stats())
	}
	return stats
}

// Count 活跃会话数
func (p *SessionPool) Count() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.sessions)
}

// Cleanup 清理死亡的会话
func (p *SessionPool) Cleanup() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	removed := 0
	for id, s := range p.sessions {
		if !s.IsAlive() {
			s.Close()
			delete(p.sessions, id)
			removed++
		}
	}
	return removed
}
