package v1

import (
	"github.com/kamva/mgm/v3"
)

type TaskGroup struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name"`
}

func NewTaskGroup(name string) *TaskGroup {
	return &TaskGroup {
		Name:        name,
	}
}