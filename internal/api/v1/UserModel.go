package v1

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	CreateDate       time.Time `json:"createDate"`
	Email            string    `json:"email"`
	Password         string    `json:"password"`
}
