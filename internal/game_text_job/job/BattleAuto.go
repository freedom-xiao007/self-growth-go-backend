package job

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"seltGrowth/internal/api/v1/game_text_auto"
	"time"
)

type BattleAuto interface {
	Battle(username string) (message string, isWin bool, hero game_text_auto.Hero, enemy game_text_auto.Enemy, err error)
}

type battleAuto struct {
}

func NewBattleAuto() BattleAuto {
	return &battleAuto{
	}
}

func (b *battleAuto) Battle(username string) (message string, isWin bool, hero game_text_auto.Hero, enemy game_text_auto.Enemy, err error) {
	var user game_text_auto.GameUser
	err = mgm.Coll(&user).First(bson.M{"username": username}, &user)
	if err != nil {
		return "数据库异常", true, hero, enemy, err
	}

	battleHeroes := make([]game_text_auto.Hero, 0)
	for _, hero := range user.OwnHero {
		if hero.IsBattle {
			battleHeroes = append(battleHeroes, hero)
		}
	}

	if len(battleHeroes) < 1 {
		return "无上阵英雄，失败", false, hero, enemy, nil
	}

	rand.Seed(time.Now().Unix())
	roundIndex := rand.Int63n(int64(len(battleHeroes)))
	hero = battleHeroes[roundIndex]

	enemy, err = game_text_auto.RoundGenerateEnemy()
	if err != nil {
		return "随机怪物生成异常", true, hero, enemy, err
	}

	heroDamage := hero.Strong - enemy.Defence
	enemyDamage := enemy.Strong - hero.Defence
	msg := fmt.Sprintf(": hero attck = %d, defence = %d, enemy attack = %d, defence = %d, heroDamage=%d, enenmyDamage=%d",
		hero.Strong, hero.Defence, enemy.Strong, enemy.Defence, heroDamage, enemyDamage)
	if heroDamage < 1 && enemyDamage < 1 {
		return "旗鼓相当，怪物最终撤退" + msg, true, hero, enemy, err
	}
	if heroDamage < enemyDamage {
		return "英雄不敌，失败" + msg, false, hero, enemy, err
	}
	if heroDamage > enemyDamage {
		return "英雄神武，击杀怪物" + msg, true, hero, enemy, err
	}
	if hero.AttackSpeed > enemy.AttackSpeed {
		return "旗鼓相当，但英雄速度更胜一筹，击杀怪物" + msg, true, hero, enemy, err
	}
	if hero.AttackSpeed < enemy.AttackSpeed {
		return "旗鼓相当，但怪物速度更胜一筹，被击杀,败北" + msg, false, hero, enemy, err
	}
	return "旗鼓相当，怪物最终撤退" + msg, true, hero, enemy, err
}
