package v1

import (
	"github.com/kamva/mgm/v3"
)

type TaskConfig struct {
	mgm.DefaultModel `bson:",inline"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Label       string `json:"label"`
	CycleType   int8   `json:"cycleType"`
	UserName    string `json:"userName"`
	IsComplete  bool   `json:"isComplete"`
}

func NewTaskConfig(name, description, label string, cycleType int8) *TaskConfig {
	return &TaskConfig {
		Name: name,
		Description: description,
		Label: label,
		CycleType: cycleType,
		IsComplete: false,
	}
}

func (t *TaskConfig) RefreshStatus(records []TaskRecord) {
	//now := time.Now()
	//cycleType := config.CycleType
	t.IsComplete = false
}