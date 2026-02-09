package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Result 标准响应结构
type Result struct {
	Code int         `json:"code"` // 业务状态码，0表示成功，非0表示错误
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 数据载体
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

// SuccessWithMsg 成功响应（自定义消息）
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  msg,
		Data: data,
	})
}

// Error 错误响应
func Error(c *gin.Context, httpCode int, code int, msg string) {
	c.JSON(httpCode, Result{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// BadRequest 400 错误
func BadRequest(c *gin.Context, msg string) {
	Error(c, http.StatusBadRequest, http.StatusBadRequest, msg)
}

// InternalError 500 错误
func InternalError(c *gin.Context, msg string) {
	Error(c, http.StatusInternalServerError, http.StatusInternalServerError, msg)
}
