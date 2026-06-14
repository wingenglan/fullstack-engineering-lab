package handler

import (
	"net/http"

	"fullstack-engineering-lab/server/internal/model"
	"fullstack-engineering-lab/server/internal/response"
	"fullstack-engineering-lab/server/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeAuthFailed, "invalid request: "+err.Error())
		return
	}

	if err := h.authService.Register(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeAuthFailed, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeAuthFailed, "invalid request: "+err.Error())
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, response.CodeAuthFailed, err.Error())
		return
	}

	response.Success(c, resp)
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, response.CodeTokenExpired, "user not found in context")
		return
	}

	resp, err := h.authService.GetProfile(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusNotFound, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		response.Error(c, http.StatusBadRequest, response.CodeAuthFailed, "token not found")
		return
	}

	if err := h.authService.Logout(token.(string)); err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, "failed to logout")
		return
	}

	response.Success(c, nil)
}
