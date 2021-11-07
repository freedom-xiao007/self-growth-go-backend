package main

import (
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"seltGrowth/internal/api/v1/game_text_auto"
	"seltGrowth/internal/pkg/store/mongodb"
)

// 初始化敌人
func main() {
	enemies := []string{
		"贪欲",
		"嗔恨",
		"痴",
		"傲慢",
		"犹疑",
		"生苦",
		"老苦",
		"病苦 ",
		"死苦",
		"爱别离苦",
		"怨憎会苦",
		"所求不得苦",
		"五取蕴苦",
	}

	mongodb.InitMongodb()

	for _, hero := range enemies {
		var heroExist game_text_auto.Enemy
		err := mgm.Coll(&heroExist).First(bson.M{"name_py": hero}, &heroExist)
		if err != nil {
			newHero := game_text_auto.NewHero(hero, hero)
			err := mgm.Coll(&heroExist).Create(newHero)
			if err != nil {
				log.Error(err)
				return
			}
			log.Info("生成：：", hero)
		} else {
			log.Info("已存在：：", hero)
		}
	}
}
