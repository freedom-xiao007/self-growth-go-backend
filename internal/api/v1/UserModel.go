package v1

import (
	"github.com/kamva/mgm/v3"
	"time"
)

// User 创建时间
// LearnPersistDay 学习持续天数
// RunningPersistDay 运动持续天数
// SleepPersistDay 睡觉持续天数
// ImprovePersistDay 进步持续天数
type User struct {
	mgm.DefaultModel  `bson:",inline"`
	CreateDate        time.Time `json:"createDate"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	LearnPersistDay   int64     `json:"learnPersistDay"`
	RunningPersistDay int64     `json:"runningPersistDay"`
	SleepPersistDay   int64     `json:"sleepPersistDay"`
	ImprovePersistDay int64     `json:"improvePersistDay"`
}
