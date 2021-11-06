package achievement

import "github.com/kamva/mgm/v3"

// Achievement 每日成就
// Date 日期 天
// Flesh 肉身点 == 运动
// Spirit 灵魂点 == 睡觉
// Divinity 神性 == 学习
type Achievement struct {
	mgm.DefaultModel `bson:",inline"`
	Date string `json:"date"`
	Flesh int64 `json:"flesh"`
	Spirit int64 `json:"spirit"`
	Divinity int64 `json:"divinity"`
}

func NewAchievement(date string, flesh, spirit, divinity int64) *Achievement {
	return &Achievement{
		Date: date,
		Flesh: flesh,
		Spirit: spirit,
		Divinity: divinity,
	}
}

func NewEmptyAchievement(date string) *Achievement {
	return &Achievement{
		Date: date,
		Flesh: 0,
		Spirit: 0,
		Divinity: 0,
	}
}
