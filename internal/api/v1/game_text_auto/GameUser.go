package game_text_auto

import (
	"github.com/kamva/mgm/v3"
)

// GameUser 每日的所得的游戏资源
// Username 用户名称，目前采用邮箱为唯一标识
// 睡觉得到元气(Spirit)
//     基础10点，处于6-8小时（360-480）之间翻倍加成
//     元气：增加灵魂力、神识、神防
//     灵魂力： 睡觉处于6-8小时，每日全恢复，其他扣20%
// 运动得到精元(Strength)
//     分钟数 == 获取精元数
//     精元：提升血气、力量、命中率、攻速、闪避、防御、移动力
// 学习与完成任务，得到灵气(Reiki)
//     学习分钟数 == 获得的灵气值
//     任务换算：1输出类：非博客笔记 == 20灵气，1输出类：博客笔记 == 50灵气, 1其他 == 10灵脉
//     灵气：抽取角色、购买装备等等
// OwnHero 拥有的角色列表
type GameUser struct {
	mgm.DefaultModel `bson:",inline"`
	Spirit int64 `json:"spirit"`
	Strength int64 `json:"strength"`
	Reiki int64 `json:"reiki"`
	Username string `json:"username"`
	OwnHero  map[string]Hero `json:"ownHero"`
}

func NewGameUser(achievement DayAchievement) *GameUser {
	return &GameUser{
		Spirit: achievement.Spirit,
		Strength: achievement.Strength,
		Reiki: achievement.Reiki,
		Username: achievement.Username,
	}
}
