package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	modelV1 "seltGrowth/internal/api/v1"
	"seltGrowth/internal/growth_record/middleware"
	srvV1 "seltGrowth/internal/growth_record/service/v1"
	"strconv"
	"time"
)

type AchievementController struct {
	srv srvV1.UserService
}

func NewAchievementController() *AchievementController {
	return &AchievementController{
		srv: srvV1.NewAchievementService(),
	}
}

func (a *AchievementController) get(c *gin.Context) {
	timestamp, err := strconv.Atoi(c.Query("timestamp"))
	if err != nil {
		ErrorResponse(c, 400, errors.New("请传入有效时间").Error())
		return
	}
	data, err = a.srv.get(time.Unix(int64(timestamp), 0))
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}

	SuccessResponse(c, middleware.GenerateToken(c, user.Email))
}
