package dao

import (
	"github.com/FlowRays/ColorizeFlow/backend/model"
	"gorm.io/gorm"
)

type ImageDAO struct {
	DB *gorm.DB
}

func (dao *ImageDAO) Create(image *model.Image) error {
	err := dao.DB.Create(image).Error
	return err
}

func (dao *ImageDAO) GetByID(id uint) (*model.Image, error) {
	var image model.Image
	err := dao.DB.Where("id = ?", id).First(&image).Error
	return &image, err
}

func (dao *ImageDAO) GetByUsername(username string) ([]model.ImageWithUser, error) {
	var images []model.ImageWithUser
	query := dao.DB.Table("images").
		Select("images.id, users.username, images.created_at").
		Joins("JOIN users ON images.user_id = users.id").
		Where("type = ?", "painting")
	if username != "" {
		query = query.Where("username = ?", username)
	}
	err := query.Find(&images).Error
	return images, err
}
