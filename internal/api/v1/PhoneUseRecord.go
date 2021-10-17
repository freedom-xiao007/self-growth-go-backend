package v1

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type PhoneUseRecord struct {
	mgm.DefaultModel `bson:",inline"`
	Date     time.Time `json:"date"`
	Activity string    `json:"activity"`
}

func NewPhoneUserRecord(activity string) *PhoneUseRecord {
	return &PhoneUseRecord {
		Date: time.Now(),
		Activity: activity,
	}
}
