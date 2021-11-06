package service_game_text

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"seltGrowth/internal/api/v1/game_text_auto"
)

type HeroService interface {
	List() ([]game_text_auto.Hero, error)
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
