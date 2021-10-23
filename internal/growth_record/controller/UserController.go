package controller

import (
	"github.com/gin-gonic/gin"
	modelV1 "seltGrowth/internal/api/v1"
	"seltGrowth/internal/growth_record/middleware"
	srvV1 "seltGrowth/internal/growth_record/service/v1"
)

type UserController struct {
	srv srvV1.UserService
}

func NewUserController() *UserController {
	return &UserController{
		srv: srvV1.NewUserService(),
	}
}

func (t *UserController) Login(c *gin.Context) {
	var user modelV1.User
	err := c.BindJSON(&user)
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}

	if user.Email == "" || user.Password == "" {
		ErrorResponse(c, 400, "邮箱和密码不能为空")
		return
	}
	err = t.srv.Login(user)
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}

	SuccessResponse(c, middleware.GenerateToken(c, user.Email))
}
