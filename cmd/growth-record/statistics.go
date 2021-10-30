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

	for i := 10; i < 31; i++ {
		query := bson.M{}
		query["username"] = "1243925457@qq.com"
		query["date"] = bson.M{operator.Gte: time.Date(2021, 10, i, 6, 0, 0, 0, time.UTC)}
		query["date"] = bson.M{operator.Lte: time.Date(2021, 10, i+1, 6, 0, 0, 0, time.UTC)}

		var phoneUseRecords []modelV1.PhoneUseRecord
		err = mgm.Coll(&modelV1.PhoneUseRecord{}).SimpleFind(&phoneUseRecords, query)
		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println(i, "日期数据：", len(phoneUseRecords))
		//for _, item := range phoneUseRecords {
		//	log.Info(item)
		//}

		//time.Sleep(time.Second * 5)
	}

	log.Info("数据统计完毕")
}
