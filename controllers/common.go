package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 统一响应格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// PaginatedResponse 分页响应
type PaginatedResponse struct {
	Total int64       `json:"total"`
	Items interface{} `json:"items"`
}

// SuccessPaginatedResponse 成功分页响应
func SuccessPaginatedResponse(c *gin.Context, total int64, items interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PaginatedResponse{
			Total: total,
			Items: items,
		},
	})
}
