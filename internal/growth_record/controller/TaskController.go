package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
