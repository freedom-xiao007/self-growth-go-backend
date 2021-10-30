package v1

import (
	"errors"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
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
	AddTaskGroup(taskGroup modelV1.TaskGroup) error
	TaskListByGroup(username string) (map[string][]modelV1.TaskConfig, error)
	Overview(userName string, startTimeStamp, endTimeStamp int64) (interface{}, error)
	DeleteTaskGroup(groupName string, userName string) error
	DeleteTask(id string, userName string) error
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
		taskConfig.CycleType, taskConfig.LearnType, time.Now())
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

func (t *taskService) AddTaskGroup(taskGroup modelV1.TaskGroup) error {
	var existTaskGroup modelV1.TaskGroup
	err := mgm.Coll(&modelV1.TaskGroup{}).First(bson.M{"username": taskGroup.UserName, "name": taskGroup.Name}, &existTaskGroup)
	if err != nil {
		err := mgm.Coll(&modelV1.TaskGroup{}).Create(&taskGroup)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("任务组已存在:" + taskGroup.Name)
}

func (t *taskService) TaskListByGroup(username string) (map[string][]modelV1.TaskConfig, error) {
	var taskGroups []modelV1.TaskGroup
	err := mgm.Coll(&modelV1.TaskGroup{}).SimpleFind(&taskGroups, bson.M{"username": username})
	if err != nil {
		return nil, err
	}

	var taskConfigs []modelV1.TaskConfig
	err = mgm.Coll(&modelV1.TaskConfig{}).SimpleFind(&taskConfigs, bson.M{"username": username})
	if err != nil {
		return nil, err
	}

	result := make(map[string][]modelV1.TaskConfig)
	for _, value := range taskGroups {
		result[value.Name] = make([]modelV1.TaskConfig, 0)
	}

	for _, taskConfig := range taskConfigs {
		var records []modelV1.TaskRecord
		err = mgm.Coll(&modelV1.TaskRecord{}).SimpleFind(&records, bson.M{"username": username})
		if err != nil {
			return nil, err
		}
		taskConfig.RefreshStatus(records)

		result[taskConfig.Group] = append(result[taskConfig.Group], taskConfig)
	}

	return result, nil
}

func (t *taskService) Overview(userName string, startTimeStamp, endTimeStamp int64) (interface{}, error) {
	query := bson.M{}
	query["username"] = userName
	if startTimeStamp > 0 {
		query["date"] = bson.M{operator.Gte: time.Unix(startTimeStamp, 0)}
	}
	if endTimeStamp > 0 {
		query["date"] = bson.M{operator.Let: time.Unix(endTimeStamp, 0)}
	}

	activityStatistics, err := activityStatistics(query)
	if err != nil {
		return nil, err
	}
	taskStatistics, err := taskStatistics(query)
	if err != nil {
		return nil, err
	}

	response := make(map[string]interface{})
	response["activityStatistics"] = activityStatistics
	response["taskStatistics"] = taskStatistics
	return response, nil
}

func activityStatistics(query bson.M) (interface{}, error) {
	var applicationActivity []modelV1.ApplicationActivity
	err := mgm.Coll(&modelV1.ApplicationActivity{}).SimpleFind(&applicationActivity, query)
	if err != nil {
		return nil, err
	}

	activity2Application := make(map[string]string)
	for _, item := range applicationActivity {
		activity2Application[item.ActivityName] = item.ApplicationName
	}

	var phoneUseRecords []modelV1.PhoneUseRecord
	err = mgm.Coll(&modelV1.PhoneUseRecord{}).SimpleFind(&phoneUseRecords, query)
	if err != nil {
		return nil, err
	}

	statistics := make(map[string]float64)
	for _, item := range phoneUseRecords {
		if _, ok := activity2Application[item.Activity]; !ok {
			continue
		}
		application := activity2Application[item.Activity]
		if _, ok := statistics[application]; !ok {
			statistics[application] = 10.0
		} else {
			statistics[application] = statistics[application] + 10.0
		}
	}
	return statistics, nil
}

func taskStatistics(query bson.M) (interface{}, error) {
	var taskRecords []modelV1.TaskRecord
	err := mgm.Coll(&modelV1.TaskRecord{}).SimpleFind(&taskRecords, query)
	if err != nil {
		return nil, err
	}

	labelStatistics := make(map[string]int64)
	for _, item := range taskRecords {
		if _, ok := labelStatistics[item.Label]; !ok {
			labelStatistics[item.Label] = 1
		} else {
			labelStatistics[item.Label] = labelStatistics[item.Label] + 1
		}
	}

	statistics := make(map[string]interface{})
	statistics["labelStatistics"] = labelStatistics
	statistics["amount"] = len(taskRecords)
	return statistics, nil
}

func (t *taskService) DeleteTaskGroup(groupName string, userName string) error {
	var taskGroup modelV1.TaskGroup
	err := mgm.Coll(&modelV1.TaskGroup{}).First(bson.M{"name": groupName, "username": userName}, &taskGroup)
	if err != nil {
		return err
	}
	if taskGroup.UserName != userName {
		return errors.New("越权删除他人数据")
	}

	var taskConfigs []modelV1.TaskConfig
	err = mgm.Coll(&modelV1.TaskConfig{}).SimpleFind(&taskConfigs, bson.M{"username": userName, "group": groupName})
	if err != nil {
		return err
	}

	for _, taskConfig := range taskConfigs {
		err = mgm.Coll(&modelV1.TaskConfig{}).Delete(&taskConfig)
		if err != nil {
			return err
		}
	}


	err = mgm.Coll(&modelV1.TaskGroup{}).Delete(&taskGroup)
	if err != nil {
		return err
	}
	return nil
}

func (t *taskService) DeleteTask(id string, userName string) error {
	var taskConfig modelV1.TaskConfig
	err := mgm.Coll(&modelV1.TaskConfig{}).FindByID(id, &taskConfig)
	if err != nil {
		return err
	}
	if taskConfig.UserName != userName {
		return errors.New("越权删除他人数据")
	}
	err = mgm.Coll(&modelV1.TaskConfig{}).Delete(&taskConfig)
	if err != nil {
		return err
	}
	return nil
}