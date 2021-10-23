package v1

import (
	"errors"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	modelV1 "seltGrowth/internal/api/v1"
)

type LabelService interface {
	AddLabel(label modelV1.LabelModel) error
	LabelList(userName string) ([]modelV1.LabelModel, error)
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

func (l *labelService) LabelList(username string) ([]modelV1.LabelModel, error) {
	var labels []modelV1.LabelModel
	err := mgm.Coll(&modelV1.LabelModel{}).SimpleFind(&labels, bson.M{"username": username})
	if err != nil {
		return nil, err
	}
	return labels, nil
}