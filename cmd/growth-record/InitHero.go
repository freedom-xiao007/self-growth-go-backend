package main

import (
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"seltGrowth/internal/api/v1/game_text_auto"
	"seltGrowth/internal/pkg/store/mongodb"
)

// 初始化英雄池
func main() {
	heroes := []string{
		"东皇太一",
		"昊天上帝",
		"盘古大帝",
		"女娲娘娘",
		"太昊",
		"炎帝",
		"黄帝",
		"少昊 ",
		"颛顼",
		"火神",
	}

	mongodb.InitMongodb()

	for _, hero := range heroes {
		var heroExist game_text_auto.Hero
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
