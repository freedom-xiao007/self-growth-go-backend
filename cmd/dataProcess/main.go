package main

import (
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	modelV1 "seltGrowth/internal/api/v1"
)

func main() {
	err := mgm.SetDefaultConfig(nil, "phone_record", options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	var records []modelV1.PhoneUseRecord
	err = mgm.Coll(&modelV1.PhoneUseRecord{}).SimpleFind(&records, bson.M{})
	if err != nil {
		log.Error(err)
		return
	}

	for _, record := range records {
		err := mgm.Coll(&modelV1.PhoneUseRecord{}).Update(&record)
		if err != nil {
			log.Error(err)
		}
	}

	log.Info("数据更新完毕")
}
