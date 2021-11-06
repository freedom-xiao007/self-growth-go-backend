package service_game_text

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"seltGrowth/internal/api/v1/game_text_auto"
)

type HeroService interface {
	List() ([]game_text_auto.Hero, error)
	UserInfo(userName string) (game_text_auto.GameUser, error)
}

type heroService struct {
}

func NewHeroService() HeroService {
	return &heroService{}
}

func (h heroService) List() ([]game_text_auto.Hero, error) {
	var heroes []game_text_auto.Hero
	err := mgm.Coll(&game_text_auto.Hero{}).SimpleFind(&heroes, bson.M{})
	if err != nil {
		return nil, err
	}
	return heroes, nil
}

func (h heroService) UserInfo(userName string) (game_text_auto.GameUser, error) {
	var user game_text_auto.GameUser
	err := mgm.Coll(&user).First(bson.M{"username": userName}, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}