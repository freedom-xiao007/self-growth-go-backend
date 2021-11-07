package game_text_auto

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type BattleLog struct {
	mgm.DefaultModel `bson:",inline"`
	Date time.Time `json:"date"`
	Message string `json:"message"`
	IsWin bool `json:"isWin"`
	Hero Hero `json:"hero"`
	Enemy Enemy `json:"enemy"`
	Username string `json:"username"`
}
