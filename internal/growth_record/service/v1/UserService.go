package v1

import (
	"errors"
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/scrypt"
	modelV1 "seltGrowth/internal/api/v1"
	"time"
)

type UserService interface {
	Login(user modelV1.User) error
}

type userService struct {

}

func NewUserService() UserService{
	return &userService{}
}

func (u *userService) Login(user modelV1.User) error {
	log.Infoln(user)
	var existUser modelV1.User
	err := mgm.Coll(&modelV1.User{}).First(bson.M{"email": user.Email}, &existUser)
	if err != nil {
		err := createNewUser(user)
		if err != nil {
			return err
		}
		return nil
	}

	isPass, err := checkPassword(user, existUser)
	if err != nil {
		return err
	}
	if !isPass {
		return errors.New("密码错误")
	}
	log.Infoln("用户登录：", user.Email)
	return nil
}

func createNewUser(user modelV1.User) error {
	encrypt, err := encryptPassword(user)
	if err != nil {
		return err
	}
	user.CreateDate = time.Now()
	user.Password = encrypt
	err = mgm.Coll(&modelV1.User{}).Create(&user)
	if err != nil {
		return err
	}
	return nil
}

func checkPassword(user, existUser modelV1.User) (bool, error) {
	encrypt, err := encryptPassword(user)
	if err != nil {
		return false, err
	}
	if existUser.Password != encrypt {
		return false, nil
	}
	return true, nil
}

func encryptPassword(user modelV1.User) (string, error) {
	dk, err := scrypt.Key([]byte(user.Password), []byte(user.Email), 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}
	return string(dk), nil
}