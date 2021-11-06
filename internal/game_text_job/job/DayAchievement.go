package job

import (
	"encoding/json"
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	v1 "seltGrowth/internal/api/v1"
	"seltGrowth/internal/api/v1/game_text_auto"
	"seltGrowth/internal/pkg/service/statistics"
	"time"
)

type DayAchievementCal interface {
	Cal() error
}

type dayAchievementCal struct {
	dayService statistics.DayStatisticsService
}

func NewDayAchievementCal() DayAchievementCal {
	return &dayAchievementCal{
		dayService: statistics.NewDayStatisticsService(),
	}
}

func (d *dayAchievementCal) Cal() error {
	log.Infof("每日成就统计转换开始：%s", time.Now())

	var users []v1.User
	err := mgm.Coll(&v1.User{}).SimpleFind(&users, bson.M{})
	if err != nil {
		log.Error(err)
		return err
	}

	yesterday := time.Now().AddDate(0, 0, -1)
	for _, user := range users {
		dayStatistics, err := d.dayService.DayStatistics(yesterday, user.Email, true, false)
		if err != nil {
			log.Error(err)
			return err
		}
		achievement := game_text_auto.NewDayAchievement(dayStatistics)
		err = mgm.Coll(&game_text_auto.DayAchievement{}).Create(achievement)
		if err != nil {
			log.Error(err)
			return err
		}

		s, err := json.MarshalIndent(achievement, "", "    ")
		log.Infof("昨日成就：：%s", string(s))
	}
	return nil
}
