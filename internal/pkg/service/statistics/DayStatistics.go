package statistics

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

type DayStatisticsService interface {
	DayStatistics(day time.Time, userName string, refresh bool, showAll bool) (modelV1.DayStatistics, error)
}

type dayStatisticsService struct {
}

func NewDayStatisticsService() DayStatisticsService {
	return &dayStatisticsService{}
}

func (t *dayStatisticsService) DayStatistics(day time.Time, userName string, refresh bool, showAll bool) (modelV1.DayStatistics, error) {
	endTime := time.Date(day.Year(), day.Month(), day.Day(), 22, 0, 0, 0, day.Location())
	startTime := endTime.AddDate(0, 0, -1)
	date := fmt.Sprintf("%04d-%02d-%02d", endTime.Year(), endTime.Month(), endTime.Day())
	log.Info("day statistics:: ", date)

	var existDayStatistics modelV1.DayStatistics
	_ = mgm.Coll(&modelV1.DayStatistics{}).First(bson.M{"username": userName, "date": date}, &existDayStatistics)
	if !refresh && !existDayStatistics.IsEmpty() {
		log.Info("统计存在，直接读取记录")
		return existDayStatistics, nil
	}

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
	log.Infoln(startTime, endTime)

	activity2Application := make(map[string]modelV1.ActivityModel)
	for _, activity := range activities {
		activity2Application[activity.Activity] = activity
	}

	query := bson.M{
		"username": userName,
		"date": bson.M{operator.Gte: startTime, operator.Lte: endTime},
	}
	log.Info(query)

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"date", 1}})

	var phoneUseRecords []modelV1.PhoneUseRecord
	err = mgm.Coll(&modelV1.PhoneUseRecord{}).SimpleFind(&phoneUseRecords, query, findOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	latestActivity := ""
	latestActivityDate := time.Time{}
	activityLog := make([]modelV1.ActivityLog, 0)
	activitySet := make(map[string]bool)
	activityAmount := make(map[string]int64)
	activityDateLog := make(map[string][]time.Time)
	for _, item := range phoneUseRecords {
		activity := item.Activity
		if _, ok := activity2Application[activity]; !showAll && !ok {
			if latestActivity == "" {
				continue
			}

			if _, ok := activityAmount[latestActivity]; !ok {
				activityAmount[latestActivity] = int64(item.Date.Sub(latestActivityDate).Minutes())
				activityDateLog[latestActivity] = make([]time.Time, 0)
				activityDateLog[latestActivity] = append(activityDateLog[latestActivity], latestActivityDate.In(time.Local))
				activityDateLog[latestActivity] = append(activityDateLog[latestActivity], item.Date.In(time.Local))
				activitySet[latestActivity] = true
			} else {
				activityAmount[latestActivity] = activityAmount[latestActivity] + int64(item.Date.Sub(latestActivityDate).Minutes())
				activityDateLog[latestActivity] = append(activityDateLog[latestActivity], latestActivityDate.In(time.Local))
				activityDateLog[latestActivity] = append(activityDateLog[latestActivity], item.Date.In(time.Local))
			}

			latestActivity = ""
			continue
		}

		if latestActivity == "" {
			latestActivity = activity
			latestActivityDate = item.Date
		} else if latestActivity != activity {
			if _, ok := activityAmount[latestActivity]; !ok {
				activityAmount[latestActivity] = int64(item.Date.Sub(latestActivityDate).Minutes())
				activityDateLog[latestActivity] = make([]time.Time, 0)
				activityDateLog[latestActivity] = append(activityDateLog[latestActivity], latestActivityDate.In(time.Local))
				activityDateLog[latestActivity] = append(activityDateLog[latestActivity], item.Date.In(time.Local))
				activitySet[latestActivity] = true
			} else {
				activityAmount[latestActivity] = activityAmount[latestActivity] + int64(item.Date.Sub(latestActivityDate).Minutes())
				activityDateLog[latestActivity] = append(activityDateLog[latestActivity], latestActivityDate.In(time.Local))
				activityDateLog[latestActivity] = append(activityDateLog[latestActivity], item.Date.In(time.Local))
			}

			latestActivity = activity
			latestActivityDate = item.Date
		}
	}

	for key, _ := range activitySet {
		activityLog = append(activityLog, *modelV1.NewActivityLog(key, activity2Application[key].Application, activity2Application[key].Label, activityAmount[key], activityDateLog[key]))
	}
	return activityLog, nil
}

func getTaskStatistics(userName string, startTime, endTime time.Time) (int64, []modelV1.TaskRecord, error) {
	query := bson.M{
		"username": userName,
		"completedate": bson.M{operator.Gte: startTime, operator.Lte: endTime},
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"completeDate", 1}})

	var records []modelV1.TaskRecord
	err := mgm.Coll(&modelV1.TaskRecord{}).SimpleFind(&records, query, findOptions)
	if err != nil {
		return 0, nil, err
	}
	return int64(len(records)), records, nil
}
