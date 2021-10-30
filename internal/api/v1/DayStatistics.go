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
	ActivityLog map[string]ActivityLog `json:"activityLog"`
	UserName         string `json:"userName"`
}

func (d *DayStatistics) IsEmpty() bool {
	return reflect.DeepEqual(*d, DayStatistics{})
}

func NewDayStatistics(date string, completeTaskAmount int64, completeTaskLog []TaskRecord, activityLog map[string]ActivityLog) *DayStatistics {
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
	Amount int64 `json:"amount"`
	Dates  []time.Time `json:"dates"`
}

func NewActivityLog(name string, amount int64, dates []time.Time) *ActivityLog {
	return &ActivityLog{
		Name: name,
		Amount: amount,
		Dates: dates,
	}
}