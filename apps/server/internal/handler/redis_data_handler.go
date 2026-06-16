package handler

import (
	"net/http"
	"strconv"

	"fullstack-engineering-lab/server/internal/model"
	"fullstack-engineering-lab/server/internal/response"
	"fullstack-engineering-lab/server/internal/service"

	"github.com/gin-gonic/gin"
)

// RedisDataHandler 提供 Hash / Set / ZSet / List / String 五种数据结构的演示接口
type RedisDataHandler struct {
	svc *service.RedisDataService
}

func NewRedisDataHandler(svc *service.RedisDataService) *RedisDataHandler {
	return &RedisDataHandler{svc: svc}
}

// ===== Hash 操作 =====

// SetField 设置单个 Hash 字段
func (h *RedisDataHandler) SetField(c *gin.Context) {
	var req model.HashFieldUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.svc.SetHashField(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetProfile 查询完整 Hash 画像
func (h *RedisDataHandler) GetProfile(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "缺少 key 参数")
		return
	}

	resp, err := h.svc.GetHashProfile(key)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// MultiSet 批量设置 Hash 字段
func (h *RedisDataHandler) MultiSet(c *gin.Context) {
	var req model.HashMultiSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.svc.MultiSetHash(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// DeleteField 删除 Hash 字段
func (h *RedisDataHandler) DeleteField(c *gin.Context) {
	var req struct {
		Key   string `json:"key" binding:"required"`
		Field string `json:"field" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.svc.DeleteHashField(req.Key, req.Field)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// ===== Set 操作 =====

// SetAdd 添加集合成员
func (h *RedisDataHandler) SetAdd(c *gin.Context) {
	var req model.SetAddMembersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.svc.SetAddMembers(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// SetRemove 移除集合成员
func (h *RedisDataHandler) SetRemove(c *gin.Context) {
	var req model.SetRemoveMembersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.svc.SetRemoveMembers(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// SetMembers 查询集合全部成员
func (h *RedisDataHandler) SetMembers(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "缺少 key 参数")
		return
	}

	resp, err := h.svc.GetSetMembers(&model.SetListRequest{Key: key})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// SetIntersect 计算交集
func (h *RedisDataHandler) SetIntersect(c *gin.Context) {
	var req model.SetOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.svc.SetIntersect(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// SetUnion 计算并集
func (h *RedisDataHandler) SetUnion(c *gin.Context) {
	var req model.SetOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.svc.SetUnion(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// SetDiff 计算差集
func (h *RedisDataHandler) SetDiff(c *gin.Context) {
	var req model.SetOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.svc.SetDiff(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// ===== ZSet 操作 =====

// ZSetAddScore 添加/增加成员分数
func (h *RedisDataHandler) ZSetAddScore(c *gin.Context) {
	var req model.ZSetAddScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.svc.ZSetAddScore(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// ZSetTopN 获取前 N 名
func (h *RedisDataHandler) ZSetTopN(c *gin.Context) {
	var req model.ZSetTopNRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.svc.ZSetGetTopN(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// ZSetRank 查询某成员排名
func (h *RedisDataHandler) ZSetRank(c *gin.Context) {
	var req model.ZSetRankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.svc.ZSetGetRank(&req)
	if err != nil {
		response.Error(c, http.StatusNotFound, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, resp)
}

// ===== List 操作 =====

// ListPush 推入消息
func (h *RedisDataHandler) ListPush(c *gin.Context) {
	var req model.ListPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.svc.ListPush(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// ListPop 弹出消息
func (h *RedisDataHandler) ListPop(c *gin.Context) {
	var req model.ListPopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.svc.ListPop(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// ListRange 查询列表范围
func (h *RedisDataHandler) ListRange(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "缺少 key 参数")
		return
	}

	startStr := c.DefaultQuery("start", "0")
	stopStr := c.DefaultQuery("stop", "19")
	start, _ := strconv.ParseInt(startStr, 10, 64)
	stop, _ := strconv.ParseInt(stopStr, 10, 64)

	resp, err := h.svc.ListRange(&model.ListRangeRequest{Key: key, Start: start, Stop: stop})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// ===== String 操作 =====

// StringSet 设置键值
func (h *RedisDataHandler) StringSet(c *gin.Context) {
	var req model.StringSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.svc.StringSet(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// StringGet 获取键值
func (h *RedisDataHandler) StringGet(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "缺少 key 参数")
		return
	}

	resp, err := h.svc.StringGet(key)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// StringIncr 自增计数
func (h *RedisDataHandler) StringIncr(c *gin.Context) {
	var req model.StringIncrRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.svc.StringIncr(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}
