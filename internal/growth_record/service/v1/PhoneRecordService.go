package v1

import (
	"database/sql"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	modelV1 "seltGrowth/internal/api/v1"
)

type PhoneRecordService interface {
	AddRecord(record *modelV1.PhoneUseRecord) error
	Overview(username string) ([]map[string]interface{}, error)
	ActivityHistory(activityName string, startTime, endTime sql.NullTime) ([]modelV1.PhoneUseRecord, error)
}

type phoneRecordService struct {

}

func NewPhoneRecordService() PhoneRecordService {
	return &phoneRecordService{}
}

func (p *phoneRecordService) AddRecord(record *modelV1.PhoneUseRecord) error {
	err := mgm.Coll(record).Create(record)
	if err != nil {
		return err
	}
	return nil
}

func (p *phoneRecordService) Overview(username string) ([]map[string]interface{}, error) {
	var records []modelV1.PhoneUseRecord
	err := mgm.Coll(&modelV1.PhoneUseRecord{}).SimpleFind(&records, bson.M{})
	if err != nil {
		return nil, err
	}

	statistics := make(map[string]int64)
	for _, value := range records {
		if statistics[value.Activity] < 1 {
			statistics[value.Activity] = 1
		} else {
			statistics[value.Activity] = statistics[value.Activity] + 1
		}
	}

	activity2Application, err := getActivity2Application(username)
	if err != nil {
		return nil, err
	}

	activityList := make([]map[string]interface{}, 0)
	for key, value := range statistics {
		activity := make(map[string]interface{})
		activity["name"] = key
		activity["application"] = activity2Application[key]
		activity["times"] = value
		activityList = append(activityList, activity)
	}

	return activityList, nil
}

func getActivity2Application(username string) (map[string]string, error) {
	var activities []modelV1.ActivityModel
	err := mgm.Coll(&modelV1.ActivityModel{}).SimpleFind(&activities, bson.M{"username": username})
	if err != nil {
		return nil, err
	}

	activity2Application := make(map[string]string)
	for _, activity := range activities {
		activity2Application[activity.Activity] = activity.Application
	}
	return activity2Application, nil
}

func (p *phoneRecordService) ActivityHistory(activityName string, startTime, endTime sql.NullTime) ([]modelV1.PhoneUseRecord, error) {
	var records []modelV1.PhoneUseRecord
	query := bson.M{}
	query["activity"] = activityName
	if startTime.Valid {
		query["date"] = bson.M{operator.Gte: startTime.Time}
	}
	if startTime.Valid {
		query["date"] = bson.M{operator.Let: endTime.Time}
	}
	err := mgm.Coll(&modelV1.PhoneUseRecord{}).SimpleFind(&records, query)
	if err != nil {
		return nil, err
	}
	return records, nil
}