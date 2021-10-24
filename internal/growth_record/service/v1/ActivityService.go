package v1

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	v1 "seltGrowth/internal/api/v1"
)

type ActivityService interface {
	GetActivities(username string) ([]v1.ActivityModel, error)
	UpdateActivityName(model v1.ActivityModel) (v1.ActivityModel, error)
}

type activityService struct {

}

func NewActivityService() ActivityService {
	return &activityService{}
}

func (a *activityService) GetActivities(username string) ([]v1.ActivityModel, error) {
	var records []v1.PhoneUseRecord
	mgm.Coll(&v1.PhoneUseRecord{}).SimpleFind(&records, bson.M{})

	allActivity := make(map[string]bool)
	for _, value := range records {
		if _, ok := allActivity[value.Activity]; ok {
			continue
		}
		allActivity[value.Activity] = true
	}

	var activities []v1.ActivityModel
	mgm.Coll(&v1.ActivityModel{}).SimpleFind(&activities, bson.M{})
	nameActivity := make(map[string]bool)
	for _, value := range activities {
		if _, ok := nameActivity[value.Activity]; ok {
			continue
		}
		nameActivity[value.Activity] = true
	}

	//res := make([]v1.ActivityModel, 0)
	for k := range allActivity {
		if _, ok := nameActivity[k]; ok {
			continue
		}
		activity := v1.NewActivityModel("", k, username)
		activities = append(activities, *activity)
	}
	return activities, nil
}

func (a *activityService) UpdateActivityName(model v1.ActivityModel) (v1.ActivityModel, error) {
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
	return v1.ActivityModel{}, nil
}