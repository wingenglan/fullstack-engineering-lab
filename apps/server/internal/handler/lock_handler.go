package handler

import (
	"net/http"

	"fullstack-engineering-lab/server/internal/model"
	"fullstack-engineering-lab/server/internal/response"
	"fullstack-engineering-lab/server/internal/service"

	"github.com/gin-gonic/gin"
)

type LockHandler struct {
	lockService *service.LockService
}

func NewLockHandler(lockService *service.LockService) *LockHandler {
	return &LockHandler{lockService: lockService}
}

// Acquire 获取分布式锁
// @Summary     获取分布式锁
// @Description 尝试获取指定资源的分布式锁
// @Tags        Redis Lock
// @Accept      json
// @Produce     json
// @Param       body body     model.LockAcquireRequest true "锁请求"
// @Success     200  {object} response.Response
// @Failure     409  {object} response.Response
// @Router      /lock/acquire [post]
func (h *LockHandler) Acquire(c *gin.Context) {
	var req model.LockAcquireRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeLockConflict, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.lockService.AcquireLock(&req)
	if err != nil {
		response.Error(c, http.StatusConflict, response.CodeLockConflict, err.Error())
		return
	}

	response.Success(c, resp)
}

// Release 释放分布式锁
// @Summary     释放分布式锁
// @Description 释放指定资源的分布式锁
// @Tags        Redis Lock
// @Accept      json
// @Produce     json
// @Param       body body     model.LockReleaseRequest true "释放请求"
// @Success     200  {object} response.Response
// @Failure     400  {object} response.Response
// @Router      /lock/release [post]
func (h *LockHandler) Release(c *gin.Context) {
	var req model.LockReleaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeLockConflict, "请求参数无效: "+err.Error())
		return
	}

	if err := h.lockService.ReleaseLock(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeLockConflict, err.Error())
		return
	}

	response.SuccessWithMessage(c, "锁已释放", nil)
}

// Status 查询锁状态
// @Summary     查询锁状态
// @Description 查询指定资源的分布式锁状态
// @Tags        Redis Lock
// @Accept      json
// @Produce     json
// @Param       body body     model.LockStatusRequest true "查询请求"
// @Success     200  {object} response.Response{data=model.LockStatusResponse}
// @Router      /lock/status [post]
func (h *LockHandler) Status(c *gin.Context) {
	var req model.LockStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeLockConflict, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.lockService.GetStatus(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// ContentionDemo 并发争抢演示
// @Summary     并发争抢演示
// @Description 模拟多个协程并发争抢同一把分布式锁
// @Tags        Redis Lock
// @Accept      json
// @Produce     json
// @Param       body body     model.ContentionDemoRequest true "争抢参数"
// @Success     200  {object} response.Response{data=model.ContentionDemoResponse}
// @Router      /lock/contention [post]
func (h *LockHandler) ContentionDemo(c *gin.Context) {
	var req model.ContentionDemoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeLockConflict, "请求参数无效: "+err.Error())
		return
	}

	resp, err := h.lockService.ContentionDemo(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}
