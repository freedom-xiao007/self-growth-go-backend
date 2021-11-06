package v1

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"seltGrowth/internal/api/v1/game_text_auto"
	"time"
)

type AchievementService interface {
	Get(endTime time.Time, username string) ([]game_text_auto.DayAchievement, error)
}

type achievementService struct {
}

func NewAchievementService() AchievementService {
	return &achievementService{}
}

func (a *achievementService) Get(endTime time.Time, username string) ([]game_text_auto.DayAchievement, error) {
	startTime := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 22, 0, 0, 0, endTime.Location())
	startTime = startTime.AddDate(0, 0, -10)
	query := bson.M{
		"username": username,
		"date": bson.M{operator.Gte: startTime, operator.Lte: endTime},
	}
	log.Info(query)

	var achievements []game_text_auto.DayAchievement
	err := mgm.Coll(&game_text_auto.DayAchievement{}).SimpleFind(&achievements, query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	dayAchievements := make(map[string]game_text_auto.DayAchievement)
	for _, item := range achievements {
		dayAchievements[item.Date] = item
	}

	result := make([]game_text_auto.DayAchievement, 0)
	for i:=0; i < 10; i++ {
		day := endTime.AddDate(0, 0, -i)
		date := fmt.Sprintf("%04d-%02d-%02d", day.Year(), day.Month(), day.Day())
		if _, ok := dayAchievements[date]; ok {
			result = append(result, dayAchievements[date])
		} else {
			result = append(result, *game_text_auto.NewEmptyDayAchievement(date))
		}
	}
	return result, nil
}