package v1

import (
	"time"
)

type AchievementService interface {
	Get(unix time.Time) (interface{}, error)
}

type achievementService struct {
}

func NewAchievementService() AchievementService {
	return &achievementService{}
}

func (a *achievementService) Get(unix time.Time) (interface{}, error) {
	return nil, nil
}