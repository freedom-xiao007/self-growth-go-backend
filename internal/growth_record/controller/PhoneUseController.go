package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	v1 "seltGrowth/internal/api/v1"
	srvV1 "seltGrowth/internal/growth_record/service/v1"
)

type PhoneUseController struct {
	srv srvV1.PhoneRecordService
}

func NewPhoneUseController() *PhoneUseController {
	return &PhoneUseController{
		srv: srvV1.NewPhoneRecordService(),
	}
}

func (p *PhoneUseController) UploadRecord(c *gin.Context) {
	activity := c.PostForm("activity")
	record := v1.NewPhoneUserRecord(activity)
	err := p.srv.AddRecord(record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "upload success")
}

