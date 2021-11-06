package v1

import (
	"errors"
	"fmt"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"seltGrowth/internal/api/v1/game_text_auto"
	"seltGrowth/internal/pkg/service/statistics"
	"time"
)

type AchievementService interface {
	Get(endTime time.Time, username string) ([]game_text_auto.DayAchievement, error)
	Sync(current time.Time, name string) error
}

type achievementService struct {
	dayStatistics statistics.DayStatisticsService
}

func NewAchievementService() AchievementService {
	return &achievementService{
		dayStatistics: statistics.NewDayStatisticsService(),
	}
}

func (a *achievementService) Get(endTime time.Time, username string) ([]game_text_auto.DayAchievement, error) {
	startTime := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 22, 0, 0, 0, endTime.Location())
	startTime = startTime.AddDate(0, 0, -10)
	query := bson.M{
		"username": username,
		"created_at": bson.M{operator.Gte: startTime, operator.Lte: endTime},
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
			result = append(result, *game_text_auto.NewEmptyDayAchievement(date, username))
		}
	}
	return result, nil
}

func (a *achievementService) Sync(current time.Time, username string) error {
	date := fmt.Sprintf("%04d-%02d-%02d", current.Year(), current.Month(), current.Day())
	var achievementExist game_text_auto.DayAchievement
	err := mgm.Coll(&game_text_auto.DayAchievement{}).First(bson.M{"username": username, "date": date}, &achievementExist)
	if err == nil {
		return errors.New("不可重复同步")
	}

	dayStatistics, err := a.dayStatistics.DayStatistics(current, username, true, false)
	if err != nil {
		return err
	}

	achievement := game_text_auto.NewDayAchievement(dayStatistics)
	err = mgm.Coll(&game_text_auto.DayAchievement{}).Create(achievement)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}