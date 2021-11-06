package service_game_text

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"seltGrowth/internal/api/v1/game_text_auto"
	"time"
)

type HeroService interface {
	List() ([]game_text_auto.Hero, error)
	UserInfo(userName string) (game_text_auto.GameUser, error)
	HeroRound(userName string) (string, error)
	OwnHeroes(userName string) ([]game_text_auto.Hero, error)
}

type heroService struct {
}

func NewHeroService() HeroService {
	return &heroService{}
}

func (h *heroService) List() ([]game_text_auto.Hero, error) {
	var heroes []game_text_auto.Hero
	err := mgm.Coll(&game_text_auto.Hero{}).SimpleFind(&heroes, bson.M{})
	if err != nil {
		return nil, err
	}
	return heroes, nil
}

func (h *heroService) UserInfo(userName string) (game_text_auto.GameUser, error) {
	var user game_text_auto.GameUser
	err := mgm.Coll(&user).First(bson.M{"username": userName}, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (h *heroService) HeroRound(userName string) (string, error) {
	var user game_text_auto.GameUser
	err := mgm.Coll(&user).First(bson.M{"username": userName}, &user)
	if err != nil {
		return "", err
	}
	if user.Reiki < 100 {
		return "碎片不足1百", nil
	}

	var heroes []game_text_auto.Hero
	err = mgm.Coll(&game_text_auto.Hero{}).SimpleFind(&heroes, bson.M{})
	if err != nil {
		return "", err
	}

	rand.Seed(time.Now().Unix())
	roundIndex := rand.Int63n(int64(len(heroes)))
	hero := heroes[roundIndex]

	if user.OwnHero == nil {
		user.OwnHero = make(map[string]game_text_auto.Hero)
	}
	ownHeroes := user.OwnHero
	if _, ok := ownHeroes[hero.NamePY]; ok {
		heroInfo := ownHeroes[hero.NamePY]
		heroInfo.Chip = heroInfo.Chip + 10
		user.OwnHero[hero.NamePY] = heroInfo
		user.Reiki = user.Reiki - 100
		err = mgm.Coll(&user).Update(&user)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("得到英雄：%s, 已拥有，得到10碎片", hero.NameZW), nil
	}

	user.OwnHero[hero.NamePY] = *game_text_auto.NewHero(hero.NameZW, hero.NameZW)
	user.Reiki = user.Reiki - 100
	err = mgm.Coll(&user).Update(&user)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("获得新英雄：%s", hero.NameZW), nil
}

func (h *heroService) OwnHeroes(userName string) (heroes []game_text_auto.Hero, err error) {
	var user game_text_auto.GameUser
	err = mgm.Coll(&user).First(bson.M{"username": userName}, &user)
	if err != nil {
		return nil, err
	}

	if user.OwnHero == nil {
		return heroes, nil
	}

	for _, value := range user.OwnHero {
		heroes = append(heroes, value)
	}
	return heroes, nil
}
