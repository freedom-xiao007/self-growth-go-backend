package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	v1 "seltGrowth/internal/api/v1"
	srvV1 "seltGrowth/internal/growth_record/service/v1"
	"strconv"
	"time"
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
	log.Println("phone use record:" + activity)
	record := v1.NewPhoneUserRecord(activity, c.GetHeader("userName"))
	err := p.srv.AddRecord(record)
	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}
	SuccessResponse(c, "upload success")
}

func (p *PhoneUseController) Overview(c *gin.Context) {
	data, err := p.srv.Overview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, data)
}

func (p *PhoneUseController) ActivityHistory(c *gin.Context) {
	activityName := c.Query("activity")
	start, err := strconv.Atoi(c.Query("startTimeStamp"))
	end, err := strconv.Atoi(c.Query("endTimeStamp"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	startTime := sql.NullTime{Valid: false}
	endTime := sql.NullTime{Valid: false}
	if start > 0 {
		startTime = sql.NullTime{Valid: true, Time: time.Unix(int64(start), 0)}
	}
	if end > 0 {
		endTime = sql.NullTime{Valid: true, Time: time.Unix(int64(end), 0)}
	}
	data, err := p.srv.ActivityHistory(activityName, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, data)
}

