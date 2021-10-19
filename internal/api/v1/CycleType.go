package v1

import (
	"github.com/kamva/mgm/v3"
)

type CycleType struct {
	mgm.DefaultModel `bson:",inline"`
	Name        string  `json:"activity"`
	UserName    string  `json:"userName"`
	Description string  `json:"description"`
	Flag        int8    `json:"flag"`
}

func NewCycleType(name, description string, flag int8) *CycleType {
	return &CycleType {
		Name: name,
		Description: description,
		Flag: flag,
	}
}
