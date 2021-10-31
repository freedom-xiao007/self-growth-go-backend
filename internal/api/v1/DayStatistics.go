package v1

import (
	"github.com/kamva/mgm/v3"
	"reflect"
	"time"
)

type DayStatistics struct {
	mgm.DefaultModel `bson:",inline"`
	Date string `json:"date"`
	CompleteTaskAmount int64 `json:"completeTaskAmount"`
	CompleteTaskLog    []TaskRecord `json:"completeTaskLog"`
	ActivityLog []ActivityLog `json:"activityLog"`
	UserName         string `json:"userName"`
}

func (d *DayStatistics) IsEmpty() bool {
	return reflect.DeepEqual(*d, DayStatistics{})
}

func NewDayStatistics(date string, completeTaskAmount int64, completeTaskLog []TaskRecord, activityLog []ActivityLog) *DayStatistics {
	return &DayStatistics{
		Date: date,
		CompleteTaskAmount: completeTaskAmount,
		CompleteTaskLog: completeTaskLog,
		ActivityLog: activityLog,
	}
}

type ActivityLog struct {
	Name string `json:"name"`
	Application string `json:"application"`
	Amount float64 `json:"amount"`
	Dates  []time.Time `json:"dates"`
	Label string `json:"label"`
}

func NewActivityLog(name, application, label string, amount int64, dates []time.Time) *ActivityLog {
	return &ActivityLog{
		Name: name,
		Amount: float64(amount) * 10 / 60,
		Dates: dates,
		Application: application,
		Label: label,
	}
}