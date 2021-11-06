package job

import (
	"seltGrowth/internal/pkg/store/mongodb"
	"testing"
)

func TestNewDayAchievementCal(t *testing.T) {
	mongodb.InitMongodb()
	o := NewDayAchievementCal()
	err := o.Cal()
	if err != nil {
		t.Error(err)
	}
}
