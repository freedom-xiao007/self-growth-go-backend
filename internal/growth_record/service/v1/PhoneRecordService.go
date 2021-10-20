package v1

import (
	"database/sql"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	v1 "seltGrowth/internal/api/v1"
)

type PhoneRecordService interface {
	AddRecord(record *v1.PhoneUseRecord) error
	Overview() (map[string]int64, error)
	ActivityHistory(activityName string, startTime, endTime sql.NullTime) ([]v1.PhoneUseRecord, error)
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

func (p *phoneRecordService) ActivityHistory(activityName string, startTime, endTime sql.NullTime) ([]v1.PhoneUseRecord, error) {
	var records []v1.PhoneUseRecord
	query := bson.M{}
	query["activity"] = activityName
	if startTime.Valid {
		query["date"] = bson.M{operator.Gte: startTime.Time}
	}
	if startTime.Valid {
		query["date"] = bson.M{operator.Let: endTime.Time}
	}
	err := mgm.Coll(&v1.PhoneUseRecord{}).SimpleFind(&records, query)
	if err != nil {
		return nil, err
	}
	return records, nil
}