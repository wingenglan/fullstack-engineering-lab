package handler

import (
	"fmt"
	"fullstack-engineering-lab/server/internal/model"
	"fullstack-engineering-lab/server/internal/response"
	"fullstack-engineering-lab/server/internal/service"
	tcpPkg "fullstack-engineering-lab/server/pkg/tcp"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TCPHandler TCP HTTP 处理器
type TCPHandler struct {
	svc *service.TCPService
}

// NewTCPHandler 创建 TCP Handler
func NewTCPHandler(svc *service.TCPService) *TCPHandler {
	return &TCPHandler{svc: svc}
}

// --- 基础 TCP 协议 ---

// CreateSession 创建 TCP 会话
func (h *TCPHandler) CreateSession(c *gin.Context) {
	sess, err := h.svc.CreateSession()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}
	response.Success(c, sess)
}

// CloseSession 关闭会话
func (h *TCPHandler) CloseSession(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.CloseSession(id); err != nil {
		response.Error(c, http.StatusNotFound, 40400, err.Error())
		return
	}
	response.Success(c, gin.H{"session_id": id, "closed": true})
}

// ListSessions 列出所有会话
func (h *TCPHandler) ListSessions(c *gin.Context) {
	stats := h.svc.ListSessions()
	response.Success(c, stats)
}

// SendCommand 发送命令
func (h *TCPHandler) SendCommand(c *gin.Context) {
	id := c.Param("id")

	var req model.TCPCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	result, err := h.svc.SendCommand(id, req.Command)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}
	response.Success(c, result)
}

// StreamSession 会话事件 SSE 流
func (h *TCPHandler) StreamSession(c *gin.Context) {
	id := c.Param("id")

	eventCh, err := h.svc.StreamSession(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, 40400, err.Error())
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	c.Stream(func(w io.Writer) bool {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case event, ok := <-eventCh:
				if !ok {
					c.SSEvent("disconnected", tcpPkg.SessionEvent{
						Type:      "disconnected",
						SessionID: id,
						Timestamp: time.Now().Format(time.RFC3339),
					})
					return false
				}
				c.SSEvent("session_event", event)
			case <-ticker.C:
				c.SSEvent("heartbeat", tcpPkg.SessionEvent{
					Type:      "heartbeat",
					SessionID: id,
					Timestamp: time.Now().Format(time.RFC3339),
				})
			case <-c.Request.Context().Done():
				return false
			}
		}
	})
}

// GetStats 获取服务器统计
func (h *TCPHandler) GetStats(c *gin.Context) {
	stats := h.svc.GetStats()
	response.Success(c, stats)
}

// --- 聊天室 ---

// CreateChatSession 创建聊天会话
func (h *TCPHandler) CreateChatSession(c *gin.Context) {
	var req model.TCPChatSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	sess, err := h.svc.CreateChatSession(req.Nickname, req.Room)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}
	response.Success(c, gin.H{
		"session_id": sess.SessionID,
		"nickname":   req.Nickname,
		"room":       req.Room,
		"created_at": sess.CreatedAt,
	})
}

// SendChatMessage 发送聊天消息
func (h *TCPHandler) SendChatMessage(c *gin.Context) {
	id := c.Param("id")

	var req model.TCPChatMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	if err := h.svc.SendChatMessage(id, req.Content); err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}
	response.Success(c, gin.H{
		"session_id": id,
		"sent":       true,
		"timestamp":  time.Now().Format(time.RFC3339),
	})
}

// ListChatRooms 列出聊天室
func (h *TCPHandler) ListChatRooms(c *gin.Context) {
	rooms := h.svc.ListChatRooms()
	response.Success(c, rooms)
}

// GetChatRoomInfo 获取聊天室详情
func (h *TCPHandler) GetChatRoomInfo(c *gin.Context) {
	name := c.Param("name")
	info := h.svc.GetChatRoomInfo(name)
	if info == nil {
		response.Error(c, http.StatusNotFound, 40400, "聊天室不存在")
		return
	}
	response.Success(c, info)
}

// GetChatRoomMessages 获取聊天室历史消息
func (h *TCPHandler) GetChatRoomMessages(c *gin.Context) {
	name := c.Param("name")
	limit := 50
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	messages := h.svc.GetChatRoomMessages(name, limit)
	response.Success(c, messages)
}

// ChatStreamSession 聊天会话 SSE 流
func (h *TCPHandler) ChatStreamSession(c *gin.Context) {
	id := c.Param("id")

	eventCh, err := h.svc.StreamSession(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, 40400, err.Error())
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	c.Stream(func(w io.Writer) bool {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case event, ok := <-eventCh:
				if !ok {
					c.SSEvent("disconnected", tcpPkg.SessionEvent{
						Type:      "disconnected",
						SessionID: id,
						Timestamp: time.Now().Format(time.RFC3339),
					})
					return false
				}
				c.SSEvent("chat_event", event)
			case <-ticker.C:
				c.SSEvent("heartbeat", tcpPkg.SessionEvent{
					Type:      "heartbeat",
					SessionID: id,
					Timestamp: time.Now().Format(time.RFC3339),
				})
			case <-c.Request.Context().Done():
				return false
			}
		}
	})
}
