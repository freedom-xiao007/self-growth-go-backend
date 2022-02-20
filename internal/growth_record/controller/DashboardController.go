package controller

import (
	"github.com/gin-gonic/gin"
	srvV1 "seltGrowth/internal/growth_record/service/v1"
)

type DashboardController struct {
	srv srvV1.DashboardService
}

func NewDashboardController() *DashboardController {
	return &DashboardController{
		srv: srvV1.NewDashboardService(),
	}
}

func (d *DashboardController) Statistics(c *gin.Context) {
	data, err := d.srv.Statistics(GetLoginUserName(c))
	if err != nil {
		ErrorResponse(c, 500, err.Error())
	}
	SuccessResponse(c, data)
}
