package v1

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	modelV1 "seltGrowth/internal/api/v1"
	"time"
)

type TaskService interface {
	GetTaskList(isComplete, username string) ([]modelV1.TaskConfig, error)
	Complete(id, username string) error
	AddTask(task modelV1.TaskConfig) error
}

type taskService struct {

}

func NewTaskService() TaskService {
	return &taskService{}
}

func (t *taskService) GetTaskList(isComplete, username string) ([]modelV1.TaskConfig, error) {
	var taskConfigs []modelV1.TaskConfig
	err := mgm.Coll(&modelV1.TaskConfig{}).SimpleFind(&taskConfigs, bson.M{"username": username})
	if err != nil {
		return nil, err
	}

	for _, taskConfig := range taskConfigs {
		var records []modelV1.TaskRecord
		err = mgm.Coll(&modelV1.TaskRecord{}).SimpleFind(&records, bson.M{"username": username})
		if err != nil {
			return nil, err
		}
		taskConfig.RefreshStatus(records)
	}

	return taskConfigs, nil
}

func (t *taskService) Complete(id, username string) error {
	var taskConfig modelV1.TaskConfig
	err := mgm.Coll(&modelV1.TaskConfig{}).FindByID(id, &taskConfig)
	if err != nil {
		return err
	}

	taskRecord := modelV1.NewTaskRecord(taskConfig.Name, taskConfig.Description, taskConfig.Label, username,
		taskConfig.CycleType, taskConfig.Type, time.Now())
	err = mgm.Coll(&modelV1.TaskRecord{}).Create(taskRecord)
	if err != nil {
		return err
	}
	return nil
}

func (t *taskService) AddTask(task modelV1.TaskConfig) error {
	return mgm.Coll(&task).Create(&task)
}