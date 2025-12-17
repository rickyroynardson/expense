package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func RespondJSON(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, Response{
		Message: message,
		Data:    data,
	})
}
