package controller

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/FlowRays/ColorizeFlow/backend/service"
)

type ImageController struct {
	ImageService *service.ImageService
}

func NewImageController(imageService *service.ImageService) *ImageController {
	return &ImageController{
		ImageService: imageService,
	}
}

func (controller *ImageController) CreateImage(c *gin.Context) {
	userID_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID not found"})
		return
	}
	userID, ok := userID_.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID type error"})
		return
	}
	imageType := c.PostForm("image_type")
	if imageType != "avatar" && imageType != "history" && imageType != "painting" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invaild image type"})
		return
	}
	imageData := c.PostForm("image")
	imageDataSplit := strings.Split(imageData, ",")
	if len(imageDataSplit) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Base64 image"})
		return
	}
	imageFormat := ""
	if strings.Contains(imageDataSplit[0], "jpeg") {
		imageFormat = ".jpg"
	} else if strings.Contains(imageDataSplit[0], "png") {
		imageFormat = ".png"
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported image format"})
		return
	}
	imageBytes, err := base64.StdEncoding.DecodeString(imageDataSplit[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Base64 image"})
		return
	}
	err = controller.ImageService.Create(c, imageFormat, imageBytes, userID, imageType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully!"})
}

func (controller *ImageController) GetByUsername(c *gin.Context) {
	username := c.PostForm("username")
	images, err := controller.ImageService.GetByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Images query error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Image queried successfully!", "images": images})
}

func (controller *ImageController) GetByID(c *gin.Context) {
	id_ := c.Param("id")
	id, err := strconv.ParseUint(id_, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ImageID error"})
		return
	}

	imageBytes, err := controller.ImageService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image file"})
		return
	}
	c.Header("Content-Type", "image/jpeg")
	c.Data(http.StatusOK, "image/jpeg", imageBytes)
}
