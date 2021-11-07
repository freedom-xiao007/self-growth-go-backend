package job

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"seltGrowth/internal/api/v1/game_text_auto"
	"time"
)

type AwardRound interface {
	Award(gameUser game_text_auto.GameUser) string
}

type awardRound struct {
}

func NewAwardRound() AwardRound {
	return &awardRound{}
}

func (a awardRound) Award(gameUser game_text_auto.GameUser) string {
	rand.Seed(time.Now().Unix())
	if rand.Int63n(100) == 1 {
		heroes := make([]string, 0)
		for key, _ := range gameUser.OwnHero {
			heroes = append(heroes, key)
		}

		rand.Seed(time.Now().Unix())
		heroIndex := rand.Int63n(int64(len(heroes)))
		heroName := heroes[heroIndex]

		hero := gameUser.OwnHero[heroName]
		hero.Chip += 1
		gameUser.OwnHero[heroName] = hero

		err := mgm.Coll(&gameUser).Update(&gameUser)
		if err != nil {
			return ""
		}
		msg := fmt.Sprintf("%s 恭喜掉落碎片1个： %s", gameUser.Username, hero.NameZW)
		log.Info(msg)
		return msg
	}
	return ""
}
