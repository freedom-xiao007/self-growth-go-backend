package main

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	modelV1 "seltGrowth/internal/api/v1"
	"time"
)

func main() {
	username := os.Getenv("mongo_user")
	passowrd := os.Getenv("mongo_password")
	host := os.Getenv("mongo_host")
	port := os.Getenv("mongo_port")
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, passowrd, host, port)
	log.Info("mongoURI", mongoURI)
	err := mgm.SetDefaultConfig(nil, "phone_record", options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
		return
	}

	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", "2021-10-17 06:00:00", time.Local)
	if err != nil {
		log.Fatal(err)
		return
	}

	duration, err := time.ParseDuration("24h")
	if err != nil {
		log.Fatal(err)
		return
	}

	nextTime := startTime.Add(duration)
	log.Info("start:", startTime)
	log.Info("next:", nextTime)
	log.Info("now:", time.Now())

	for i := 1; i < 31; i++ {
		query := bson.M{}
		query["username"] = "1243925457@qq.com"
		query["date"] = bson.M{operator.Gte: time.Date(2021, 10, i, 6, 0, 0, 0, time.UTC)}
		query["date"] = bson.M{operator.Lte: time.Date(2021, 10, i+1, 6, 0, 0, 0, time.UTC)}
		date := fmt.Sprintf("2021-10-%02d", i)
		log.Info(date)

		activityLog, err := activityStatistics(query, date)
		if err != nil {
			log.Fatal(err)
			return
		}

		completeTaskAmount, completeTaskLog, err := taskStatistics(query, date)

		dayStatistics := *modelV1.NewDayStatistics(date, completeTaskAmount, completeTaskLog, activityLog)
		err = mgm.Coll(&modelV1.DayStatistics{}).Create(&dayStatistics)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	log.Info("数据统计完毕")
}

func taskStatistics(query bson.M, s string) (int64, []modelV1.TaskRecord, error) {
	var records []modelV1.TaskRecord
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"completeDate", 1}})
	err := mgm.Coll(&modelV1.TaskRecord{}).SimpleFind(&records, query, findOptions)
	if err != nil {
		return 0, nil, err
	}
	return int64(len(records)), records, nil
}

func activityStatistics(query bson.M, date string) (map[string]modelV1.ActivityLog, error) {
	var phoneUseRecords []modelV1.PhoneUseRecord
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"date", 1}})
	err := mgm.Coll(&modelV1.PhoneUseRecord{}).SimpleFind(&phoneUseRecords, query, findOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println(date, "日期数据：", len(phoneUseRecords))

	activityLog := make(map[string]modelV1.ActivityLog)
	activitySet := make(map[string]bool)
	activityAmount := make(map[string]int64)
	activityDateLog := make(map[string][]time.Time)
	for _, item := range phoneUseRecords {
		activity := item.Activity
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
		activityLog[key] = *modelV1.NewActivityLog(key, activityAmount[key], activityDateLog[key])
	}
	return activityLog, nil
}
