package v1

import (
	"errors"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	modelV1 "seltGrowth/internal/api/v1"
	"time"
)

type TaskService interface {
	GetTaskList(isComplete, username string) ([]modelV1.TaskConfig, error)
	Complete(id, username string) error
	AddTask(task modelV1.TaskConfig) error
	History(userName string) ([]modelV1.TaskRecord, error)
	AddTaskGroup(taskGroupName, username string) error
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

func (t *taskService) History(userName string) ([]modelV1.TaskRecord, error) {
	var records []modelV1.TaskRecord
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"completeDate", -1}})
	findOptions.SetSkip(0)
	findOptions.SetLimit(100)
	err := mgm.Coll(&modelV1.TaskRecord{}).SimpleFind(&records, bson.M{"username": userName}, findOptions)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (t *taskService) AddTaskGroup(taskGroupName, username string) error {
	var existTaskGroup modelV1.TaskGroup
	err := mgm.Coll(&modelV1.TaskGroup{}).First(bson.M{"username": username, "name": taskGroupName}, &existTaskGroup)
	if err != nil {
		err := mgm.Coll(&modelV1.TaskGroup{}).Create(modelV1.NewTaskGroup(taskGroupName, username))
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("任务组已存在:" + taskGroupName)
}
