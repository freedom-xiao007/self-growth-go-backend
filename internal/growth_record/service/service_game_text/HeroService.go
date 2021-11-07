package service_game_text

import (
	"errors"
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"seltGrowth/internal/api/v1/game_text_auto"
	"time"
)

type HeroService interface {
	List() ([]game_text_auto.Hero, error)
	UserInfo(userName string) (game_text_auto.GameUser, error)
	HeroRound(userName string) (string, error)
	OwnHeroes(userName string) ([]game_text_auto.Hero, error)
	ModifyOwnHeroProperty(heroName string, property string, modifyType string, userName string) error
	BattleHero(heroName string, userName string) error
	BattleLog(userName string, pageIndex int, pageSize int) (logs []game_text_auto.BattleLog, total int64, err error)
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

// ModifyOwnHeroProperty
// todo 需要重构优化代码;属性不能减到负的处理
func (h *heroService) ModifyOwnHeroProperty(heroName string, property string, modifyType string, userName string) error {
	var user game_text_auto.GameUser
	err := mgm.Coll(&user).First(bson.M{"username": userName}, &user)
	if err != nil {
		return err
	}
	hero := user.OwnHero[heroName]
	modifyValue := 1
	if modifyType == "-1" {
		modifyValue = -1
	}

	if property == "spirit" || property == "spiritAttack" || property == "spiritDefence" {
		return modifySpiritProperty(property, int64(modifyValue), user, hero, heroName)
	}
	if property == "level" {
		return modifyLevelProperty(user, hero, heroName)
	}
	return modifyPhysicalProperty(property, int64(modifyValue), user, hero, heroName)
}

func modifySpiritProperty(property string, modifyValue int64, user game_text_auto.GameUser, hero game_text_auto.Hero, heroName string) error {
	if user.Spirit <= 0 {
		return errors.New("元气不足")
	}
	if property == "spirit" {
		hero.Spirit = hero.Spirit + modifyValue
	} else if property == "spiritAttack" {
		hero.SpiritAttack = hero.SpiritAttack + modifyValue
	} else if property == "spiritDefence" {
		hero.SpiritDefence = hero.SpiritDefence + modifyValue
	}
	user.Spirit = user.Spirit + -modifyValue

	user.OwnHero[heroName] = hero
	err := mgm.Coll(&user).Update(&user)
	if err != nil {
		return err
	}
	return nil
}

func modifyLevelProperty(user game_text_auto.GameUser, hero game_text_auto.Hero, heroName string) error {
	needChip := hero.Level * 10
	if hero.Chip < needChip {
		return errors.New("碎片不足")
	}

	hero.Level = hero.Level + 1
	hero.Chip = hero.Chip - 10
	user.OwnHero[heroName] = hero
	err := mgm.Coll(&user).Update(&user)
	if err != nil {
		return err
	}
	return nil
}

func modifyPhysicalProperty(property string, modifyValue int64, user game_text_auto.GameUser, hero game_text_auto.Hero, heroName string) error {
	if user.Strength <= 0 {
		return errors.New("精元不足")
	}
	if property == "bleed" {
		hero.Bleed = hero.Bleed + modifyValue
	} else if property == "strong" {
		hero.Strong = hero.Strong + modifyValue
	} else if property == "shooting" {
		hero.Shooting = hero.Shooting + modifyValue
	} else if property == "attackSpeed" {
		hero.AttackSpeed = hero.AttackSpeed + modifyValue
	} else if property == "dodge" {
		hero.Dodge = hero.Dodge + modifyValue
	} else if property == "defence" {
		hero.Defence = hero.Defence + modifyValue
	} else if property == "moveSpeed" {
		hero.MoveSpeed = hero.MoveSpeed + modifyValue
	}
	user.Strength = user.Strength + -modifyValue

	user.OwnHero[heroName] = hero
	err := mgm.Coll(&user).Update(&user)
	if err != nil {
		return err
	}
	return nil
}

// BattleHero todo 获取影响或许可以抽成一个通用的方法
func (h *heroService) BattleHero(heroName string, userName string) error {
	var user game_text_auto.GameUser
	err := mgm.Coll(&user).First(bson.M{"username": userName}, &user)
	if err != nil {
		return err
	}
	hero := user.OwnHero[heroName]
	hero.IsBattle = !hero.IsBattle
	user.OwnHero[heroName] = hero
	err = mgm.Coll(&user).Update(&user)
	if err != nil {
		return err
	}
	return nil
}

func (h *heroService) BattleLog(userName string, pageIndex int, pageSize int) (logs []game_text_auto.BattleLog, total int64, err error) {
	query := bson.M{"username": userName}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"date", -1}})
	findOptions.SetSkip(int64(pageIndex * pageSize))
	findOptions.SetLimit(int64(pageSize))
	err = mgm.Coll(&game_text_auto.BattleLog{}).SimpleFind(&logs, query, findOptions)
	if err != nil {
		return nil, 0, err
	}

	total, err = mgm.Coll(&game_text_auto.BattleLog{}).CountDocuments(mgm.Ctx(), query)
	if err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}