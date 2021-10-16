package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	srvV1 "seltGrowth/internal/growth_record/service/v1"
)

type HelloHandler struct {
	srv srvV1.Service
}

func NewHelloHandler() *HelloHandler {
	return &HelloHandler{
		srv: srvV1.NewService(),
	}
}

func (hello *HelloHandler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello")
}
