package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorResponse(c *gin.Context, code int16, message string) {
	c.JSON(int(code), gin.H{
		"code":    code,
		"message": message,
	})
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": data,
	})
}

func GetLoginUserName(c *gin.Context) string {
	return c.GetHeader("userName")
}
