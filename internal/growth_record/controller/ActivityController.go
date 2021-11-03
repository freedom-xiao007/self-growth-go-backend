package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	modelV1 "seltGrowth/internal/api/v1"
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

func (a *ActivityController) UploadRecord(c *gin.Context) {
	activity := c.PostForm("activity")
	log.Println("phone use record:" + activity)
	record := modelV1.NewPhoneUserRecord(activity, GetLoginUserName(c))
	err := a.srv.AddRecord(record)
	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}
	SuccessResponse(c, "upload success")
}

func (a *ActivityController) Overview(c *gin.Context) {
	data, err := a.srv.Overview(GetLoginUserName(c))
	if err != nil {
		log.Error(err)
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, data)
}

func (a *ActivityController) ActivityHistory(c *gin.Context) {
	activityName := c.Query("activity")

	start, err := strconv.Atoi(c.Query("startTimeStamp"))
	end, err := strconv.Atoi(c.Query("endTimeStamp"))
	startTime := sql.NullTime{Valid: false}
	endTime := sql.NullTime{Valid: false}
	if start > 0 {
		startTime = sql.NullTime{Valid: true, Time: time.Unix(int64(start), 0)}
	}
	if end > 0 {
		endTime = sql.NullTime{Valid: true, Time: time.Unix(int64(end), 0)}
	}

	size := 100
	if c.Query("pageSize") != "" {
		size, err = strconv.Atoi(c.Query("pageSize"))
		if err != nil {
			ErrorResponse(c, 400, err.Error())
			return
		}
	}
	index := 0
	if c.Query("pageIndex") != "" {
		index, err = strconv.Atoi(c.Query("pageIndex"))
		if err != nil {
			ErrorResponse(c, 400, err.Error())
			return
		}
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	data, total, err := a.srv.ActivityHistory(GetLoginUserName(c), activityName, startTime, endTime, index, size)
	if err != nil {
		log.Error(err)
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, map[string]interface{}{"data": data, "total": total})
}

func (a *ActivityController) UpdateActivityModel(c *gin.Context) {
	var activityModel modelV1.ActivityModel
	err := c.BindJSON(&activityModel)
	if err != nil {
		log.Error(err)
		ErrorResponse(c, 400, err.Error())
		return
	}

	activityModel.UserName = GetLoginUserName(c)
	err = a.srv.UpdateActivityModel(activityModel)
	if err != nil {
		log.Error(err)
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, "更新成功")
}
