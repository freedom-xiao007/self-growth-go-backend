package v1

import (
	"github.com/kamva/mgm/v3"
)

type LabelModel struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	UserName         string `json:"userName"`
}

func NewLabelModel(name, description, username string) *LabelModel {
	return &LabelModel{
		Name:        name,
		Description: description,
		UserName:    username,
	}
}
