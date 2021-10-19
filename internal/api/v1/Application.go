package v1

import (
	"github.com/kamva/mgm/v3"
)

type Application struct {
	mgm.DefaultModel `bson:",inline"`
	Name     string  `json:"name"`
	UserName string  `json:"userName"`
}

func NewApplication(name string) *Application {
	return &Application {
		Name: name,
	}
}

type ApplicationActivity struct {
	mgm.DefaultModel `bson:",inline"`
	ApplicationName  string    `json:"applicationName"`
	ActivityName     string    `json:"activityName"`
	UserName         string    `json:"userName"`
}

func NewApplicationActivity(applicationName, activityName string) *ApplicationActivity {
	return &ApplicationActivity {
		ApplicationName: applicationName,
		ActivityName: activityName,
	}
}