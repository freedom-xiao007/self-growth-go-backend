package main

import (
	"github.com/kamva/mgm/v3"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	v1 "seltGrowth/internal/api/v1"
	"seltGrowth/internal/api/v1/game_text_auto"
	"seltGrowth/internal/game_text_job/job"
	"seltGrowth/internal/pkg/store/mongodb"
	"time"
)

func main() {
	mongodb.InitMongodb()

	log.Info("自我生长-文字版游戏")
	log.Info("定时任务初始化")

	dayAchievementCal := job.NewDayAchievementCal()
	battleAuto := job.NewBattleAuto()
	winAward := job.NewAwardRound()

	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("0 10 6 * * *", func() {
		err := dayAchievementCal.Cal()
		if err != nil {
			log.Error(err)
		}
	})

	_, err = c.AddFunc("0/5 * * * * *", func() {
		var users []v1.User
		err := mgm.Coll(&v1.User{}).SimpleFind(&users, bson.M{})
		if err != nil {
			log.Error(err)
			return
		}

		for _, user := range users {
			log.Infof("%s :: 自动战斗开始", user.Email)
			msg, win, hero, enemy, err := battleAuto.Battle(user.Email)
			if err != nil {
				log.Error(err)
				continue
			}

			// 奖励结算
			if win {
				var gameUser game_text_auto.GameUser
				err = mgm.Coll(&gameUser).First(bson.M{"username": user.Email}, &gameUser)
				if err != nil {
					log.Error(err)
					continue
				}
				msg += "\n" + winAward.Award(gameUser)
			}

			battleLog := game_text_auto.BattleLog{
				Date: time.Now(),
				Message: msg,
				IsWin: win,
				Hero: hero,
				Enemy: enemy,
				Username: user.Email,
			}
			err = mgm.Coll(&battleLog).Create(&battleLog)
			if err != nil {
				log.Error(err)
				return
			}
			log.Info("战斗结果：:", battleLog)
		}
	})
	if err != nil {
		log.Error(err)
		return
	}

	c.Start()

	select {}
}
