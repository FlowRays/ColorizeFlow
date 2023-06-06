package dao

import (
	"github.com/FlowRays/ColorizeFlow/backend/model"
	"gorm.io/gorm"
)

type UserDAO struct {
	DB *gorm.DB
}

func (dao *UserDAO) Create(user *model.User) error {
	err := dao.DB.Create(user).Error
	return err
}

func (dao *UserDAO) GetByID(id int) (*model.User, error) {
	var user model.User
	err := dao.DB.Where("id = ?", user.ID).First(&user).Error
	return &user, err
}

func (dao *UserDAO) GetByName(username string) (*model.User, error) {
	var user model.User
	err := dao.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}
