package controller

import (
	"github.com/gin-gonic/gin"
	modelV1 "seltGrowth/internal/api/v1"
	srvV1 "seltGrowth/internal/growth_record/service/v1"
)

type LabelController struct {
	srv srvV1.LabelService
}

func NewLabelController() *LabelController {
	return &LabelController{
		srv: srvV1.NewLabelService(),
	}
}

func (l *LabelController) Add(c *gin.Context) {
	labelName := c.PostForm("name")
	if labelName == "" {
		ErrorResponse(c, 400, "名称不能为空")
		return
	}

	err := l.srv.AddLabel(*modelV1.NewLabelModel(labelName, labelName, GetLoginUserName(c)))
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, "新增成功")
}

func (l *LabelController) List(c *gin.Context) {
	data, err := l.srv.LabelList(GetLoginUserName(c))
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, data)
}
