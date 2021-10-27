package controller

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
	isComplete := c.Query("isComplete")
	tasks, err := t.srv.GetTaskList(isComplete, GetLoginUserName(c))
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, tasks)
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

func (t *TaskController) Complete(c *gin.Context) {
	id := c.PostForm("id")
	err := t.srv.Complete(id[1:len(id)-1], GetLoginUserName(c))
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, "任务完成")
}

func (t *TaskController) History(c *gin.Context) {
	data, err := t.srv.History(GetLoginUserName(c))
	if err != nil {
		log.Error(err)
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, data)
}

func (t *TaskController) AddTaskGroup(c *gin.Context) {
	taskGroupName := c.PostForm("taskGroup")
	err := t.srv.AddTaskGroup(taskGroupName, GetLoginUserName(c))
	if err != nil {
		log.Error(err)
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponseWithoutData(c)
}

func (t *TaskController) TaskListByGroup(c *gin.Context) {
	data, err := t.srv.TaskListByGroup(GetLoginUserName(c))
	if err != nil {
		log.Error(err)
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, data)
}

func (t *TaskController) Overview(c *gin.Context) {
	data, err := t.srv.Overview(GetLoginUserName(c))
	if err != nil {
		log.Error(err)
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, data)
}
