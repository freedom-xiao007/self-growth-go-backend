package v1

import (
	"fmt"
	"seltGrowth/internal/pkg/store/mongodb"
	"testing"
)

func TestDashboardService_Statistics(t *testing.T) {
	mongodb.InitMongodb()
	service := NewDashboardService()
	data, err := service.Statistics("1243925457@qq.com")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(data)
}
