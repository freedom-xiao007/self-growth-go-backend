package v1

import (
	"github.com/kamva/mgm/v3"
)

type TaskGroup struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name"`
	UserName         string `json:"username"`
}

func NewTaskGroup(name, username string) *TaskGroup {
	return &TaskGroup {
		Name:        name,
		UserName: username,
	}
}