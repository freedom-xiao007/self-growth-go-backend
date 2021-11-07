package job

import (
	"seltGrowth/internal/pkg/store/mongodb"
	"testing"
)

func TestBattleAuto_Battle(t *testing.T) {
	mongodb.InitMongodb()
	battle := NewBattleAuto()
	msg, win, hero, enemy, err := battle.Battle("1243925457@qq.com")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(msg)
	t.Log(win)
	t.Log(hero)
	t.Log(enemy)
}
