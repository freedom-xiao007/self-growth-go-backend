package v1

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	v1 "seltGrowth/internal/api/v1"
)

type PhoneRecordService interface {
	AddRecord(record *v1.PhoneUseRecord) error
	Overview() (map[string]int64, error)
}

type phoneRecordService struct {

}

func NewPhoneRecordService() PhoneRecordService {
	return &phoneRecordService{}
}

func (p *phoneRecordService) AddRecord(record *v1.PhoneUseRecord) error {
	err := mgm.Coll(record).Create(record)
	if err != nil {
		return err
	}
	return nil
}

func (p *phoneRecordService) Overview() (map[string]int64, error) {
	var records []v1.PhoneUseRecord
	mgm.Coll(&v1.PhoneUseRecord{}).SimpleFind(&records, bson.M{})

	statistics := make(map[string]int64)
	for _, value := range records {
		if statistics[value.Activity] < 1 {
			statistics[value.Activity] = 1
		} else {
			statistics[value.Activity] = statistics[value.Activity] + 1
		}
	}

	for key, value := range statistics {
		println(key, value)
	}
	return statistics, nil
}