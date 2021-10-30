package v1

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type TaskConfig struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Label            string `json:"label"`
	CycleType        int8   `json:"cycleType"`
	UserName         string `json:"userName"`
	IsComplete       bool   `json:"isComplete"`
	LearnType        int8   `json:"learnType"`
	Group            string `json:"group"`
	OutputType       int8   `json:"outputType"`
}

func NewTaskConfig(name, description, label string, cycleType int8) *TaskConfig {
	return &TaskConfig{
		Name:        name,
		Description: description,
		Label:       label,
		CycleType:   cycleType,
		IsComplete:  false,
	}
}

func (t *TaskConfig) RefreshStatus(records []TaskRecord) {
	switch t.CycleType {
	case 0:
		t.IsComplete = dailyCheck(records)
		break
	case 1:
		t.IsComplete = weekCheck(records)
		break
	case 2:
		t.IsComplete = monthCheck(records)
		break
	case 3:
		t.IsComplete = yearCheck(records)
		break
	default:
		t.IsComplete = oneTimeCheck(records)
	}
}

func dailyCheck(records []TaskRecord) bool {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, now.Location())
	for _, record := range records {
		if record.CompleteDate.After(today) {
			return true
		}
	}
	return false
}

func weekCheck(records []TaskRecord) bool {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStartDate := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	for _, record := range records {
		if record.CompleteDate.After(weekStartDate) {
			return true
		}
	}
	return false
}

func monthCheck(records []TaskRecord) bool {
	now := time.Now()
	monthStartDate := time.Date(now.Year(), now.Month(), 1, 6, 0, 0, 0, now.Location())
	for _, record := range records {
		if record.CompleteDate.After(monthStartDate) {
			return true
		}
	}
	return false
}

func yearCheck(records []TaskRecord) bool {
	now := time.Now()
	yearStartDate := time.Date(now.Year(), 1, 1, 6, 0, 0, 0, now.Location())
	for _, record := range records {
		if record.CompleteDate.After(yearStartDate) {
			return true
		}
	}
	return false
}

func oneTimeCheck(records []TaskRecord) bool {
	return len(records) > 0
}
