package v1

import (
	"database/sql"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	modelV1 "seltGrowth/internal/api/v1"
	"sort"
)

type ActivityService interface {
	GetActivities(username string) ([]modelV1.ActivityModel, error)
	UpdateActivityName(model modelV1.ActivityModel) (modelV1.ActivityModel, error)
	AddRecord(record *modelV1.PhoneUseRecord) error
	Overview(username string) ([]map[string]interface{}, error)
	ActivityHistory(username, activityName string, startTime, endTime sql.NullTime) ([]modelV1.PhoneUseRecord, error)
	UpdateActivityModel(model modelV1.ActivityModel) error
}

type activityService struct {
}

func NewActivityService() ActivityService {
	return &activityService{}
}

func (a *activityService) GetActivities(username string) ([]modelV1.ActivityModel, error) {
	var records []modelV1.PhoneUseRecord
	mgm.Coll(&modelV1.PhoneUseRecord{}).SimpleFind(&records, bson.M{})

	allActivity := make(map[string]bool)
	for _, value := range records {
		if _, ok := allActivity[value.Activity]; ok {
			continue
		}
		allActivity[value.Activity] = true
	}

	var activities []modelV1.ActivityModel
	mgm.Coll(&modelV1.ActivityModel{}).SimpleFind(&activities, bson.M{})
	nameActivity := make(map[string]bool)
	for _, value := range activities {
		if _, ok := nameActivity[value.Activity]; ok {
			continue
		}
		nameActivity[value.Activity] = true
	}

	//res := make([]modelV1.ActivityModel, 0)
	for k := range allActivity {
		if _, ok := nameActivity[k]; ok {
			continue
		}
		activity := modelV1.NewActivityModel("", k, username)
		activities = append(activities, *activity)
	}
	return activities, nil
}

func (a *activityService) UpdateActivityName(model modelV1.ActivityModel) (modelV1.ActivityModel, error) {
	var records []modelV1.PhoneUseRecord
	mgm.Coll(&modelV1.PhoneUseRecord{}).SimpleFind(&records, bson.M{})

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
	return modelV1.ActivityModel{}, nil
}

func (a *activityService) AddRecord(record *modelV1.PhoneUseRecord) error {
	err := mgm.Coll(record).Create(record)
	if err != nil {
		return err
	}
	return nil
}

func (a *activityService) Overview(username string) ([]map[string]interface{}, error) {
	var records []modelV1.PhoneUseRecord
	err := mgm.Coll(&modelV1.PhoneUseRecord{}).SimpleFind(&records, bson.M{"username": username})
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
		activity["application"] = activity2Application[key].Application
		activity["times"] = value
		activity["label"] = activity2Application[key].Label
		activityList = append(activityList, activity)
	}

	sort.Slice(activityList, func(i, j int) bool {
		return activityList[i]["times"].(int64) >= activityList[j]["times"].(int64)
	})

	return activityList, nil
}

func getActivity2Application(username string) (map[string]modelV1.ActivityModel, error) {
	var activities []modelV1.ActivityModel
	err := mgm.Coll(&modelV1.ActivityModel{}).SimpleFind(&activities, bson.M{"username": username})
	if err != nil {
		return nil, err
	}

	activity2Application := make(map[string]modelV1.ActivityModel)
	for _, activity := range activities {
		activity2Application[activity.Activity] = activity
	}
	return activity2Application, nil
}

func (a *activityService) ActivityHistory(username, activityName string, startTime, endTime sql.NullTime) ([]modelV1.PhoneUseRecord, error) {
	var records []modelV1.PhoneUseRecord
	query := bson.M{}
	query["username"] = username
	if activityName != "" {
		query["activity"] = activityName
	}
	if startTime.Valid {
		query["date"] = bson.M{operator.Gte: startTime.Time}
	}
	if startTime.Valid {
		query["date"] = bson.M{operator.Let: endTime.Time}
	}

	findOptions := options.Find()
	// Sort by `price` field descending
	findOptions.SetSort(bson.D{{"date", -1}})
	findOptions.SetSkip(0)
	findOptions.SetLimit(100)
	err := mgm.Coll(&modelV1.PhoneUseRecord{}).SimpleFind(&records, query, findOptions)
	if err != nil {
		return nil, err
	}

	activity2Application, err := getActivity2Application(username)
	if err != nil {
		return nil, err
	}

	for index, record := range records {
		records[index].Application = activity2Application[record.Activity].Application
	}
	return records, nil
}

func (a *activityService) UpdateActivityModel(activityModel modelV1.ActivityModel) error {
	var existActivityModel modelV1.ActivityModel
	err := mgm.Coll(&modelV1.ActivityModel{}).First(bson.M{"activity": activityModel.Activity, "username": activityModel.UserName}, &existActivityModel)
	if err != nil {
		err := mgm.Coll(&modelV1.ActivityModel{}).Create(&activityModel)
		if err != nil {
			return err
		}
		return nil
	}

	existActivityModel.Application = activityModel.Application
	existActivityModel.Label = activityModel.Label
	return mgm.Coll(&modelV1.ActivityModel{}).Update(&existActivityModel)
}
