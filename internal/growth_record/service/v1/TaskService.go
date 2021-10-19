package v1

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	v1 "seltGrowth/internal/api/v1"
)

type TaskService interface {
	GetTaskList() ([]v1.TaskConfig, error)
	CompleteTask(model v1.ActivityModel) (v1.ActivityModel, error)
}

type taskService struct {

}

func NewTaskService() TaskService {
	return &taskService{}
}

func (t *taskService) GetTaskList() ([]v1.TaskConfig, error) {
	var labels []v1.LabelModel
	err := mgm.Coll(&v1.LabelModel{}).SimpleFind(&labels, bson.M{})
	if err != nil {
		return nil, err
	}

	var taskConfigs []v1.TaskConfig
	err = mgm.Coll(&v1.TaskConfig{}).SimpleFind(&taskConfigs, bson.M{})
	if err != nil {
		return nil, err
	}

	for _, taskConfig := range taskConfigs {
		var records []v1.TaskRecord
		err = mgm.Coll(&v1.TaskRecord{}).SimpleFind(&records, bson.M{})
		if err != nil {
			return nil, err
		}
		taskConfig.RefreshStatus(records)
	}

	return taskConfigs, nil
}

func (t *taskService) CompleteTask(model v1.ActivityModel) (v1.ActivityModel, error) {
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