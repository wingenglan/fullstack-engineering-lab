package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 业务错误码
const (
	CodeSuccess       = 0
	CodeAuthFailed    = 40001
	CodeTokenExpired  = 40003
	CodeForbidden     = 40004
	CodeLockConflict  = 40009
	CodeInternalError = 50000
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, httpCode int, bizCode int, message string) {
	c.JSON(httpCode, Response{
		Code:    bizCode,
		Message: message,
		Data:    nil,
	})
}
