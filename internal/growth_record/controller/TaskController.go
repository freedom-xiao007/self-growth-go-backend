package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	modelV1 "seltGrowth/internal/api/v1"
	srvV1 "seltGrowth/internal/growth_record/service/v1"
)

type TaskController struct {
	srv srvV1.TaskService
}

func NewTaskController() *TaskController {
	return &TaskController{
		srv: srvV1.NewTaskService(),
	}
}

func (t *TaskController) TaskList(c *gin.Context) {
	tasks, err := t.srv.GetTaskList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (t *TaskController) AddTask(c *gin.Context) {
	var taskConfig modelV1.TaskConfig
	err := c.BindJSON(&taskConfig)
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}

	taskConfig.UserName = GetLoginUserName(c)
	err = t.srv.AddTask(taskConfig)
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, "新增任务成功")
}
