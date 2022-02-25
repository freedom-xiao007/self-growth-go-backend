package v1

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	modelV1 "seltGrowth/internal/api/v1"
	"time"
)

type DashboardService interface {
	Statistics(username string) (dashboard modelV1.DashboardStatisticsModel, err error)
}

type dashboardService struct {
}

func NewDashboardService() DashboardService {
	return &dashboardService{}
}

func (d dashboardService) Statistics(username string) (dashboard modelV1.DashboardStatisticsModel, err error) {
	day := time.Now()
	endTime := time.Date(day.Year(), day.Month(), day.Day(), 22, 0, 0, 0, day.Location())
	startTime := time.Date(day.Year(), day.Month(), day.Day()-1, 22, 0, 0, 0, day.Location())
	date := fmt.Sprintf("%04d-%02d-%02d", endTime.Year(), endTime.Month(), endTime.Day())
	log.Info("day statistics:: ", date)

	taskGroup, err := taskGroupStatistics(username, startTime, endTime)
	if err != nil {
		return dashboard, err
	}
	learnGroup, runningGroup, sleepGroup, err := activityGroupStatistics(username, startTime, endTime)
	overview := overviewGroup(learnGroup, runningGroup, sleepGroup, taskGroup)

	groups := make(map[string]modelV1.DashboardGroup)
	groups["overview"] = overview
	groups["learn"] = learnGroup
	groups["running"] = runningGroup
	groups["sleep"] = sleepGroup
	groups["task"] = taskGroup
	return *modelV1.NewDashboardStatisticsModel(groups), nil
}

func overviewGroup(learnGroup, runningGroup, sleepGroup, taskGroup modelV1.DashboardGroup) modelV1.DashboardGroup {
	timeAmount := learnGroup.Minutes + runningGroup.Minutes + sleepGroup.Minutes
	apps := make([]modelV1.DashboardApp, 0)
	apps = append(apps, *modelV1.NewDashboardApp(learnGroup.Name, learnGroup.Minutes))
	apps = append(apps, *modelV1.NewDashboardApp(runningGroup.Name, runningGroup.Minutes))
	apps = append(apps, *modelV1.NewDashboardApp(sleepGroup.Name, sleepGroup.Minutes))
	apps = append(apps, *modelV1.NewDashboardApp(taskGroup.Name, taskGroup.Minutes))
	return *modelV1.NewDashboardGroup("统计总览", timeAmount, apps)
}

func activityGroupStatistics(username string, startTime, endTime time.Time) (learnGroup, runningGroup, sleepGroup modelV1.DashboardGroup, err error) {
	var activities []modelV1.ActivityModel
	err = mgm.Coll(&modelV1.ActivityModel{}).SimpleFind(&activities, bson.M{"username": username})
	if err != nil {
		return learnGroup, runningGroup, sleepGroup, err
	}
	log.Infoln(startTime, endTime)

	activity2Application := make(map[string]modelV1.ActivityModel)
	for _, activity := range activities {
		activity2Application[activity.Activity] = activity
	}

	query := bson.M{
		"username": username,
		"date":     bson.M{operator.Gte: startTime, operator.Lte: endTime},
	}
	log.Info(query)

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"date", 1}})

	var phoneUseRecords []modelV1.PhoneUseRecord
	err = mgm.Coll(&modelV1.PhoneUseRecord{}).SimpleFind(&phoneUseRecords, query, findOptions)
	if err != nil {
		log.Fatal(err)
		return learnGroup, runningGroup, sleepGroup, err
	}

	// phoneUseRecords 表得到Activity 和 时间
	// 关联 Activity model 得到应用名和标签
	learnTimeAmount := 0
	runningTimeAmount := 0
	sleepTimeAmount := 0
	appTimeCache := make(map[string]map[string]int)
	appTimeCache["学习"] = make(map[string]int)
	appTimeCache["运动"] = make(map[string]int)
	appTimeCache["睡觉"] = make(map[string]int)
	latestActivity := ""
	latestDate := time.Now()
	endActivity := ""
	endDate := time.Now()
	for _, item := range phoneUseRecords {
		activity := item.Activity
		if _, ok := activity2Application[activity]; !ok {
			continue
		}

		endActivity = activity
		endDate = item.Date
		if latestActivity == "" {
			latestActivity = activity
			latestDate = item.Date
		} else if activity == latestActivity {
			continue
		}

		speedTime := int(item.Date.Sub(latestDate).Minutes())
		//uploadTime := item.Date
		application := activity2Application[latestActivity].Application
		label := activity2Application[latestActivity].Label

		if label == "学习" {
			learnTimeAmount = learnTimeAmount + speedTime
			appTimeCache["学习"][application] = appTimeCache["学习"][application] + speedTime
		} else if label == "运动" {
			runningTimeAmount = runningTimeAmount + speedTime
			appTimeCache["运动"][application] = appTimeCache["运动"][application] + speedTime
		} else if label == "睡觉" {
			sleepTimeAmount = sleepTimeAmount + speedTime
			appTimeCache["睡觉"][application] = appTimeCache["睡觉"][application] + speedTime
		}

		latestActivity = activity
		latestDate = item.Date
	}

	if endActivity != "" && !endDate.Equal(latestDate) {
		speedTime := int(endDate.Sub(latestDate).Minutes())
		application := activity2Application[latestActivity].Application
		label := activity2Application[latestActivity].Label
		if label == "学习" {
			learnTimeAmount = learnTimeAmount + speedTime
			appTimeCache["学习"][application] = appTimeCache["学习"][application] + speedTime
		} else if label == "运动" {
			runningTimeAmount = runningTimeAmount + speedTime
			appTimeCache["运动"][application] = appTimeCache["运动"][application] + speedTime
		} else if label == "睡觉" {
			sleepTimeAmount = sleepTimeAmount + speedTime
			appTimeCache["睡觉"][application] = appTimeCache["睡觉"][application] + speedTime
		}
	}

	learnApps := make([]modelV1.DashboardApp, 0)
	for appName, speedTime := range appTimeCache["学习"] {
		learnApps = append(learnApps, *modelV1.NewDashboardApp(appName, speedTime/60))
	}
	learnGroup.Name = "学习统计"
	learnGroup.Minutes = learnTimeAmount
	learnGroup.Apps = learnApps

	runningApps := make([]modelV1.DashboardApp, 0)
	for appName, speedTime := range appTimeCache["运动"] {
		runningApps = append(runningApps, *modelV1.NewDashboardApp(appName, speedTime/60))
	}
	runningGroup.Name = "运动统计"
	runningGroup.Minutes = runningTimeAmount
	runningGroup.Apps = runningApps

	sleepApps := make([]modelV1.DashboardApp, 0)
	for appName, speedTime := range appTimeCache["睡觉"] {
		sleepApps = append(sleepApps, *modelV1.NewDashboardApp(appName, speedTime/60))
	}
	sleepGroup.Name = "睡觉统计"
	sleepGroup.Minutes = sleepTimeAmount
	sleepGroup.Apps = sleepApps

	return learnGroup, runningGroup, sleepGroup, err
}

func taskGroupStatistics(username string, startTime, endTime time.Time) (group modelV1.DashboardGroup, err error) {
	query := bson.M{
		"username":     username,
		"completedate": bson.M{operator.Gte: startTime, operator.Lte: endTime},
	}

	var records []modelV1.TaskRecord
	mgm.Coll(&modelV1.TaskRecord{}).SimpleFind(&records, query)

	group.Name = "任务统计"
	group.Minutes = len(records)
	apps := make([]modelV1.DashboardApp, 0)
	for _, value := range records {
		apps = append(apps, *modelV1.NewDashboardApp(value.Name, 1))
	}
	group.Apps = apps
	return group, nil
}
