package v1

import (
	"github.com/kamva/mgm/v3"
)

// ActivityModel activity 手机活动名称，主键
// application 手机应用名称，能对应多个activity
type ActivityModel struct {
	mgm.DefaultModel `bson:",inline"`
	Application string  `json:"application"`
	Activity    string  `json:"activity"`
	UserName    string  `json:"userName"`
	Label       string  `json:"label"`
}

func NewActivityModel(application, activity, username string) *ActivityModel {
	return &ActivityModel {
		Application: application,
		Activity: activity,
		UserName: username,
	}
}
