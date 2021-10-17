package v1

import (
	"github.com/kamva/mgm/v3"
)

type ActivityModel struct {
	mgm.DefaultModel `bson:",inline"`
	Name     string  `json:"name"`
	Activity string  `json:"activity"`
}

func NewActivityModel(name, activity string) *ActivityModel {
	return &ActivityModel {
		Name: name,
		Activity: activity,
	}
}
