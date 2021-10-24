package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	v1 "seltGrowth/internal/api/v1"
	srvV1 "seltGrowth/internal/growth_record/service/v1"
	"strconv"
	"time"
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

func (p *ActivityController) UploadRecord(c *gin.Context) {
	activity := c.PostForm("activity")
	log.Println("phone use record:" + activity)
	record := v1.NewPhoneUserRecord(activity, GetLoginUserName(c))
	err := p.srv.AddRecord(record)
	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}
	SuccessResponse(c, "upload success")
}

func (p *ActivityController) Overview(c *gin.Context) {
	data, err := p.srv.Overview(GetLoginUserName(c))
	if err != nil {
		log.Error(err)
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, data)
}

func (p *ActivityController) ActivityHistory(c *gin.Context) {
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
