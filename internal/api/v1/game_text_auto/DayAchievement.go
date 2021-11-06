package game_text_auto

import (
	"github.com/kamva/mgm/v3"
	v1 "seltGrowth/internal/api/v1"
)

// DayAchievement 每日的所得的游戏资源
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
// IsImport 是否已导入到游戏中
type DayAchievement struct {
	mgm.DefaultModel `bson:",inline"`
	Date string `json:"date"`
	Spirit int64 `json:"spirit"`
	Strength int64 `json:"strength"`
	Reiki int64 `json:"reiki"`
	IsImport bool `json:"isImport"`
}

func NewDayAchievement(dayStatistics v1.DayStatistics) *DayAchievement {
	sleepAmount := int64(0)
	strengthAmount := int64(0)
	reikiAmount := int64(0)
	activityLogs := dayStatistics.ActivityLog
	for _, activity := range activityLogs {
		if activity.Label == "学习" {
			reikiAmount += activity.Amount
		} else if activity.Label == "运动" {
			strengthAmount += activity.Amount
		} else if activity.Label == "睡觉" {
			sleepAmount += activity.Amount
		}
	}

	taskLogs := dayStatistics.CompleteTaskLog
	for _, taskLog := range taskLogs {
		if taskLog.Type == 0 {
			reikiAmount += 10
		} else if taskLog.Type == 1 && taskLog.OutputType == 0 {
			reikiAmount += 20
		} else if taskLog.Type == 1 && taskLog.OutputType == 1 {
			reikiAmount += 50
		}
	}

	spirit := int64(10)
	if sleepAmount >= 360 && sleepAmount <= 480 {
		spirit *= 2
	}

	return &DayAchievement{
		Date: dayStatistics.Date,
		Spirit: spirit,
		Strength: strengthAmount,
		Reiki: reikiAmount,
		IsImport: false,
	}
}
