package controller

import (
	"github.com/gin-gonic/gin"
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
	activities, err := a.srv.GetActivities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, activities)
}
