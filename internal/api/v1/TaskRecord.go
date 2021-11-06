package v1

import (
	"github.com/kamva/mgm/v3"
	"time"
)

// TaskRecord 任务完成记录
// CompleteDate
// Name
// Description
// Label 任务标签，字符串：学习、睡觉、运动
// CycleType 任务周期：
// Type 任务输入输出类型：0 输入； 1 输出
// OutputType 任务输出类型：0 代码； 1 博客/笔记
type TaskRecord struct {
	mgm.DefaultModel `bson:",inline"`
	CompleteDate     time.Time `json:"completeDate"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Label            string    `json:"label"`
	CycleType        int8      `json:"cycleType"`
	UserName         string    `json:"userName"`
	Type             int8      `json:"type"`
	ConfigId         string    `json:"configId"`
	OutputType       int8      `json:"outputType"`
}

func NewTaskRecord(id, username string, completeDate time.Time, config TaskConfig) *TaskRecord {
	return &TaskRecord{
		Name: config.Name,
		Description: config.Description,
		Label: config.Label,
		CycleType: config.CycleType,
		CompleteDate: completeDate,
		Type: config.LearnType,
		UserName:     username,
		ConfigId: id,
		OutputType: config.OutputType,
	}
}
