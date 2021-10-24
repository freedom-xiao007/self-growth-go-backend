package controller

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	srvV1 "seltGrowth/internal/growth_record/service/v1"
)

type ActivityController struct {
	srv srvV1.ActivityService
}

func NewActivityController() *ActivityController {
	return &ActivityController{
		srv: srvV1.NewActivityService(),
	}
}

func (a *ActivityController) GetActivities(c *gin.Context) {
	activities, err := a.srv.GetActivities(GetLoginUserName(c))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, activities)
}
