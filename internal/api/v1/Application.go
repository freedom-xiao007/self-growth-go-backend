package v1

import (
	"github.com/kamva/mgm/v3"
)

type Application struct {
	mgm.DefaultModel `bson:",inline"`
	Name string    `json:"name"`
}

func NewApplication(name string) *Application {
	return &Application {
		Name: name,
	}
}
