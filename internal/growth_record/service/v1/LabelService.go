package v1

import (
	"errors"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	modelV1 "seltGrowth/internal/api/v1"
)

type LabelService interface {
	AddLabel(label modelV1.LabelModel) error
}

type labelService struct {

}

func NewLabelService() LabelService{
	return &labelService{}
}

func (l *labelService) AddLabel(label modelV1.LabelModel) error {
	var existLabel modelV1.LabelModel
	err := mgm.Coll(&modelV1.LabelModel{}).First(bson.M{"name": label.Name}, &existLabel)
	if err != nil {
		err := mgm.Coll(&label).Create(&label)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("已存在该标签")
}