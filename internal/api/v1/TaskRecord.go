package v1

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type TaskRecord struct {
	mgm.DefaultModel `bson:",inline"`
	CompleteDate     time.Time        `json:"date"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	Label            string           `json:"label"`
	CycleType        int8             `json:"cycleType"`
	UserName         string           `json:"userName"`
	Type             int8             `json:"type"`
}

func NewTaskRecord(name, description, label, username string, cycleType, TaskType int8, completeDate time.Time) *TaskRecord {
	return &TaskRecord {
		Name: name,
		Description: description,
		Label: label,
		CycleType: cycleType,
		CompleteDate: completeDate,
		Type: TaskType,
		UserName: username,
	}
}
