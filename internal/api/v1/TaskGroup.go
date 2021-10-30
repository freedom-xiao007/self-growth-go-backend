package v1

import (
	"github.com/kamva/mgm/v3"
)

type TaskGroup struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name"`
	UserName         string `json:"username"`
	Description      string `json:"description"`
	Label            string `json:"label"`
	CycleType        int8   `json:"cycleType"`
	IsComplete       bool   `json:"isComplete"`
	LearnType        int8   `json:"learnType"`
	ModifyId         string `json:"modifyId"`
}

func NewTaskGroup(name, username, description, label string, cycleType, learnType int8) *TaskGroup {
	return &TaskGroup {
		Name:        name,
		UserName: username,
		Description: description,
		Label: label,
		CycleType: cycleType,
		LearnType: learnType,
	}
}