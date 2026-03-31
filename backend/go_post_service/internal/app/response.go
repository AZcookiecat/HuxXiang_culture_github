package app

import "github.com/gin-gonic/gin"

type Response struct {
	Success    bool   `json:"success"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
	Pagination any    `json:"pagination,omitempty"`
	Count      *int   `json:"count,omitempty"`
}

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, Response{
		Success: true,
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, status int, message string, data any) {
	c.JSON(status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithPagination(c *gin.Context, status int, data any, pagination any) {
	c.JSON(status, Response{
		Success:    true,
		Data:       data,
		Pagination: pagination,
	})
}

func SuccessWithCount(c *gin.Context, status int, data any, count int) {
	c.JSON(status, Response{
		Success: true,
		Data:    data,
		Count:   &count,
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Success: false,
		Message: message,
	})
}

func AbortError(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}
