package mongodb

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func InitMongodb() {
	// Setup the mgm default config
	database := os.Getenv("mongo_database")
	if database == "" {
		database = "phone_record"
	}

	username := os.Getenv("mongo_user")
	password := os.Getenv("mongo_password")
	host := os.Getenv("mongo_host")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("mongo_port")
	if port == "" {
		port = "27017"
	}

	if username == "" && password == "" {
		err := mgm.SetDefaultConfig(nil, database, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	log.Info("mongoURI", mongoURI)
	err := mgm.SetDefaultConfig(nil, database, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
}
