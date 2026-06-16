package handler

import (
	"net/http"
	"strconv"

	"fullstack-engineering-lab/server/internal/model"
	"fullstack-engineering-lab/server/internal/response"
	"fullstack-engineering-lab/server/internal/service"

	"github.com/gin-gonic/gin"
)

// ChatHandler 聊天 HTTP 处理器
type ChatHandler struct {
	chatService *service.ChatService
}

// NewChatHandler 创建聊天处理器实例
func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

// WS WebSocket 连接升级
// @Summary     建立 WebSocket 连接
// @Description 将 HTTP 连接升级为 WebSocket，用于实时通信
// @Tags        Chat
// @Param       Authorization header string true "Bearer Token"
// @Success     101 {object} nil
// @Failure     401 {object} response.Response
// @Router      /chat/ws [get]
func (h *ChatHandler) WS(c *gin.Context) {
	// 从中间件获取用户信息（已在 JWTAuth 中间件设置到 context）
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, response.CodeTokenExpired, "未认证")
		return
	}

	uid := userID.(uint)

	// 用户名通过 handshake 内部查询
	handshake(uid, "", h.chatService, c)
}

// GetRooms 获取聊天室列表
// @Summary     获取聊天室列表
// @Description 获取所有可用的聊天室及其在线人数
// @Tags        Chat
// @Accept      json
// @Produce     json
// @Success     200 {object} response.Response{data=[]model.RoomResponse}
// @Router      /chat/rooms [get]
func (h *ChatHandler) GetRooms(c *gin.Context) {
	rooms, err := h.chatService.GetRooms()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, "获取聊天室列表失败")
		return
	}
	response.Success(c, rooms)
}

// CreateRoom 创建聊天室
// @Summary     创建聊天室
// @Description 创建一个新的聊天室
// @Tags        Chat
// @Accept      json
// @Produce     json
// @Param       body body     model.CreateRoomRequest true "创建请求"
// @Success     200  {object} response.Response{data=model.RoomResponse}
// @Failure     400  {object} response.Response
// @Router      /chat/rooms [post]
func (h *ChatHandler) CreateRoom(c *gin.Context) {
	var req model.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeChatError, "请求参数无效: "+err.Error())
		return
	}

	userID := c.GetUint("userID")

	room, err := h.chatService.CreateRoom(&req, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, room)
}

// GetRoomInfo 获取聊天室详情
// @Summary     获取聊天室详情
// @Description 获取指定聊天室的详细信息
// @Tags        Chat
// @Produce     json
// @Param       id  path     int true "聊天室 ID"
// @Success     200 {object} response.Response{data=model.RoomResponse}
// @Failure     404 {object} response.Response
// @Router      /chat/rooms/{id} [get]
func (h *ChatHandler) GetRoomInfo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeChatError, "无效的聊天室 ID")
		return
	}

	room, err := h.chatService.GetRoomInfo(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, response.CodeChatError, err.Error())
		return
	}

	response.Success(c, room)
}

// GetMessageHistory 获取消息历史
// @Summary     获取消息历史
// @Description 获取指定聊天室的历史消息（游标分页）
// @Tags        Chat
// @Produce     json
// @Param       room_id query    int true "聊天室 ID"
// @Param       limit   query    int false "每页数量（默认 50，最大 100）"
// @Param       before  query    int false "此 ID 之前的的消息"
// @Success     200     {object} response.Response{data=model.MessageHistoryResponse}
// @Router      /chat/messages [get]
func (h *ChatHandler) GetMessageHistory(c *gin.Context) {
	roomID, err := strconv.ParseUint(c.Query("room_id"), 10, 64)
	if err != nil || roomID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeChatError, "缺少 room_id 参数")
		return
	}

	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	var beforeID uint
	if b := c.Query("before"); b != "" {
		if parsed, err := strconv.ParseUint(b, 10, 64); err == nil {
			beforeID = uint(parsed)
		}
	}

	result, err := h.chatService.GetMessageHistory(&model.MessageHistoryRequest{
		RoomID: uint(roomID),
		Limit:  limit,
		Before: beforeID,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeChatError, err.Error())
		return
	}

	response.Success(c, result)
}

// GetOnlineUsers 获取房间在线用户
// @Summary     获取在线用户
// @Description 获取指定聊天室的在线用户列表
// @Tags        Chat
// @Produce     json
// @Param       id  path     int true "聊天室 ID"
// @Success     200 {object} response.Response
// @Router      /chat/rooms/{id}/online [get]
func (h *ChatHandler) GetOnlineUsers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeChatError, "无效的聊天室 ID")
		return
	}

	users := h.chatService.GetOnlineUsers(uint(id))
	response.Success(c, users)
}
