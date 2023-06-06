package service

import (
	"github.com/FlowRays/ColorizeFlow/backend/dao"
	"github.com/FlowRays/ColorizeFlow/backend/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImageService struct {
	ImageDAO    *dao.ImageDAO
	OSSImageDAO *dao.OSSImageDAO
	LocImageDAO *dao.LocImageDAO
}

func NewImageService(imageDAO *dao.ImageDAO, ossImageDAO *dao.OSSImageDAO, locImageDAO *dao.LocImageDAO) *ImageService {
	return &ImageService{
		ImageDAO:    imageDAO,
		OSSImageDAO: ossImageDAO,
		LocImageDAO: locImageDAO,
	}
}

func (service *ImageService) Create(c *gin.Context, imageFormat string, imageBytes []byte, userID uint, imageType string) error {
	uuid := uuid.New()
	path := imageType + "/" + uuid.String() + imageFormat
	// err := service.OSSImageDAO.Save(path, imageBytes)
	err := service.LocImageDAO.Save(path, imageBytes)
	if err != nil {
		return err
	}
	image := model.Image{Path: path, UserID: userID, Type: imageType}
	err = service.ImageDAO.Create(&image)
	return err
}

func (service *ImageService) GetByUsername(username string) ([]model.ImageWithUser, error) {
	images, err := service.ImageDAO.GetByUsername(username)
	return images, err
}

func (service *ImageService) GetByID(imageID uint) ([]byte, error) {
	image, err := service.ImageDAO.GetByID(imageID)
	if err != nil {
		return nil, err
	}
	// imageBytes, err := service.OSSImageDAO.Get(image.Path)
	imageBytes, err := service.LocImageDAO.Get(image.Path)
	if err != nil {
		return nil, err
	}
	return imageBytes, nil
}
