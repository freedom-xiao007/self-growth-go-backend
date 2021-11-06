package main

import (
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"seltGrowth/internal/game_text_job/job"
	"seltGrowth/internal/pkg/store/mongodb"
)

func main() {
	mongodb.InitMongodb()

	log.Info("自我生长-文字版游戏")
	log.Info("定时任务初始化")

	dayAchievementCal := job.NewDayAchievementCal()

	c := cron.New()
	_, err := c.AddFunc("0 10 6 * * *", func() {
		err := dayAchievementCal.Cal()
		if err != nil {
			log.Error(err)
		}
	})
	if err != nil {
		return
	}
}
