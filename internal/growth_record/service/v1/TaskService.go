package v1

import (
	"errors"
	"fmt"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	log "github.com/sirupsen/logrus"
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
	TaskListByGroup(username string) ([]map[string]interface{}, error)
	Overview(userName string, startTimeStamp, endTimeStamp int64) (interface{}, error)
	DeleteTaskGroup(groupName string, userName string) error
	DeleteTask(id string, userName string) error
	ModifyGroup(taskGroup modelV1.TaskGroup) error
	DayStatistics(day time.Time, userName string, refresh bool, showAll bool) (modelV1.DayStatistics, error)
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

	taskRecord := modelV1.NewTaskRecord(id, taskConfig.Name, taskConfig.Description, taskConfig.Label, username,
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

func (t *taskService) TaskListByGroup(username string) ([]map[string]interface{}, error) {
	var taskGroups []modelV1.TaskGroup
	err := mgm.Coll(&modelV1.TaskGroup{}).SimpleFind(&taskGroups, bson.M{"username": username})
	if err != nil {
		return nil, err
	}
	taskOfGroup := make(map[string][]modelV1.TaskConfig)
	for _, value := range taskGroups {
		taskOfGroup[value.Name] = make([]modelV1.TaskConfig, 0)
	}

	var taskConfigs []modelV1.TaskConfig
	err = mgm.Coll(&modelV1.TaskConfig{}).SimpleFind(&taskConfigs, bson.M{"username": username})
	if err != nil {
		return nil, err
	}
	for _, taskConfig := range taskConfigs {
		var records []modelV1.TaskRecord
		query := bson.M{"username": username, "configid": taskConfig.ID.Hex()}
		err = mgm.Coll(&modelV1.TaskRecord{}).SimpleFind(&records, query)
		if err != nil {
			return nil, err
		}
		taskConfig.RefreshStatus(records)
		if taskConfig.IsComplete {
			continue
		}

		taskOfGroup[taskConfig.Group] = append(taskOfGroup[taskConfig.Group], taskConfig)
	}

	result := make([]map[string]interface{}, 0)
	for _, group := range taskGroups {
		item := map[string]interface{} {"group": group, "tasks": taskOfGroup[group.Name]}
		result = append(result, item)
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
		query["date"] = bson.M{operator.Lte: time.Unix(endTimeStamp, 0)}
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

func (t *taskService) ModifyGroup(taskGroupModify modelV1.TaskGroup) error {
	var taskGroup modelV1.TaskGroup
	err := mgm.Coll(&modelV1.TaskGroup{}).FindByID(taskGroupModify.ModifyId, &taskGroup)
	if err != nil {
		return err
	}

	taskGroup.Name = taskGroupModify.Name
	return mgm.Coll(&modelV1.TaskGroup{}).Update(&taskGroup)
}


func (t *taskService) DayStatistics(day time.Time, userName string, refresh bool, showAll bool) (modelV1.DayStatistics, error) {
	startTime := time.Date(day.Year(), day.Month(), day.Day(), 6, 0, 0, 0, day.Location())
	date := fmt.Sprintf("%04d-%02d-%02d", startTime.Year(), startTime.Month(), startTime.Day())

	var existDayStatistics modelV1.DayStatistics
	_ = mgm.Coll(&modelV1.DayStatistics{}).First(bson.M{"username": userName, "date": date}, &existDayStatistics)
	if !refresh && !existDayStatistics.IsEmpty() {
		log.Info("统计存在，直接读取记录")
		return existDayStatistics, nil
	}

	endTime := startTime.AddDate(0, 0, 1)
	query := bson.M{}
	query["username"] = userName
	query["date"] = bson.M{operator.Gte: time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 6, 0, 0, 0, startTime.Location())}
	query["date"] = bson.M{operator.Lte: time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 6, 0, 0, 0, endTime.Location())}

	activityLog, err := getActivityStatistics(userName, startTime, endTime, showAll)
	if err != nil {
		return modelV1.DayStatistics{}, err
	}

	completeTaskAmount, completeTaskLog, err := getTaskStatistics(userName, startTime, endTime)

	dayStatistics := *modelV1.NewDayStatistics(date, completeTaskAmount, completeTaskLog, activityLog)
	dayStatistics.UserName = userName
	err = mgm.Coll(&modelV1.DayStatistics{}).Create(&dayStatistics)
	if err != nil {
		log.Fatal(err)
		return modelV1.DayStatistics{}, err
	}
	return dayStatistics, nil
}

func getActivityStatistics(userName string, startTime, endTime time.Time, showAll bool) ([]modelV1.ActivityLog, error) {
	var activities []modelV1.ActivityModel
	err := mgm.Coll(&modelV1.ActivityModel{}).SimpleFind(&activities, bson.M{"username": userName})
	if err != nil {
		return nil, err
	}

	activity2Application := make(map[string]modelV1.ActivityModel)
	for _, activity := range activities {
		activity2Application[activity.Activity] = activity
	}

	query := bson.M{}
	query["username"] = userName
	query["date"] = bson.M{operator.Gte: time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 6, 0, 0, 0, startTime.Location())}
	query["date"] = bson.M{operator.Lte: time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 6, 0, 0, 0, endTime.Location())}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"date", 1}})

	var phoneUseRecords []modelV1.PhoneUseRecord
	err = mgm.Coll(&modelV1.PhoneUseRecord{}).SimpleFind(&phoneUseRecords, query, findOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	activityLog := make([]modelV1.ActivityLog, 0)
	activitySet := make(map[string]bool)
	activityAmount := make(map[string]int64)
	activityDateLog := make(map[string][]time.Time)
	for _, item := range phoneUseRecords {
		activity := item.Activity
		if _, ok := activity2Application[activity]; !showAll && !ok {
			continue
		}

		if _, ok := activityAmount[activity]; !ok {
			activityAmount[activity] = 1
			activityDateLog[activity] = make([]time.Time, 0)
			activityDateLog[activity] = append(activityDateLog[activity], item.Date)
			activitySet[activity] = true
		} else {
			activityAmount[activity] = activityAmount[activity] + 1
			activityDateLog[activity] = append(activityDateLog[activity], item.Date)
		}
	}

	for key, _ := range activitySet {
		activityLog = append(activityLog, *modelV1.NewActivityLog(key, activity2Application[key].Application, activity2Application[key].Label, activityAmount[key], activityDateLog[key]))
	}
	return activityLog, nil
}

func getTaskStatistics(userName string, startTime, endTime time.Time) (int64, []modelV1.TaskRecord, error) {
	query := bson.M{}
	query["username"] = userName
	query["completedate"] = bson.M{operator.Gte: time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 6, 0, 0, 0, startTime.Location())}
	query["completedate"] = bson.M{operator.Lte: time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 6, 0, 0, 0, endTime.Location())}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"completeDate", 1}})

	var records []modelV1.TaskRecord
	err := mgm.Coll(&modelV1.TaskRecord{}).SimpleFind(&records, query, findOptions)
	if err != nil {
		return 0, nil, err
	}
	return int64(len(records)), records, nil
}