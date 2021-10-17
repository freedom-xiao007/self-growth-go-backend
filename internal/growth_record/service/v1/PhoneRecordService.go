package v1

import (
	"github.com/kamva/mgm/v3"
	v1 "seltGrowth/internal/api/v1"
)

type PhoneRecordService interface {
	AddRecord(record *v1.PhoneUseRecord) error
}

type phoneRecordService struct {

}

func NewPhoneRecordService() PhoneRecordService {
	return &phoneRecordService{}
}

func (p *phoneRecordService) AddRecord(record *v1.PhoneUseRecord) error {
	err := mgm.Coll(record).Create(record)
	if err != nil {
		return err
	}
	return nil
}
