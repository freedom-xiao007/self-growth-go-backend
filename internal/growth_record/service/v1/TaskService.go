package v1

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	v1 "seltGrowth/internal/api/v1"
)

type TaskService interface {
	GetTaskList(isComplete, username string) ([]v1.TaskConfig, error)
	CompleteTask(model v1.ActivityModel) (v1.ActivityModel, error)
	AddTask(task v1.TaskConfig) error
}

type taskService struct {

}

func NewTaskService() TaskService {
	return &taskService{}
}

func (t *taskService) GetTaskList(isComplete, username string) ([]v1.TaskConfig, error) {
	var taskConfigs []v1.TaskConfig
	err := mgm.Coll(&v1.TaskConfig{}).SimpleFind(&taskConfigs, bson.M{"username": username})
	if err != nil {
		return nil, err
	}

	for _, taskConfig := range taskConfigs {
		var records []v1.TaskRecord
		err = mgm.Coll(&v1.TaskRecord{}).SimpleFind(&records, bson.M{"username": username})
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

func (t *taskService) AddTask(task v1.TaskConfig) error {
	return mgm.Coll(&task).Create(&task)
}