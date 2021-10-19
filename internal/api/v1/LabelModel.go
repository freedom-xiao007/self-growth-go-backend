package v1

import (
	"github.com/kamva/mgm/v3"
)

type LabelModel struct {
	mgm.DefaultModel `bson:",inline"`
	Name     string  `json:"activity"`
	UserName string  `json:"userName"`
}

func NewLabelModel(name string) *LabelModel {
	return &LabelModel {
		Name: name,
	}
}
