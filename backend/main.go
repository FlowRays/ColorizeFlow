package main

import (
	"fmt"
	"log"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/FlowRays/ColorizeFlow/backend/controller"
	"github.com/FlowRays/ColorizeFlow/backend/dao"
	"github.com/FlowRays/ColorizeFlow/backend/model"
	"github.com/FlowRays/ColorizeFlow/backend/service"
)

func main() {
	err := readConfig()
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	// 连接数据库
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 执行数据库迁移
	err = db.AutoMigrate(&model.User{}, &model.Image{})
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	registerAPI(r, db)
	r.Run(":8080")
}

func readConfig() error {
	// 设置配置文件的名称和路径
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	// 读取配置文件
	err := viper.ReadInConfig()
	return err
}

func connectDB() (*gorm.DB, error) {
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUsername := viper.GetString("database.username")
	dbPassword := viper.GetString("database.password")
	dbName := viper.GetString("database.dbname")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", dbUsername, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func registerAPI(r *gin.Engine, db *gorm.DB) {
	// 依赖注入
	userDAO := &dao.UserDAO{DB: db}
	userService := service.NewUserService(userDAO)
	userController := controller.NewUserController(userService)
	imageDAO := &dao.ImageDAO{DB: db}

	ossClient, err := oss.New(
		viper.GetString("oss.endpoint"),
		viper.GetString("oss.access_key_id"),
		viper.GetString("oss.access_key_secret"),
	)
	if err != nil {
		log.Fatal(err)
	}
	ossImageDAO := &dao.OSSImageDAO{
		BucketName: viper.GetString("oss.bucket_name"),
		Client:     ossClient,
	}
	locImageDAO := &dao.LocImageDAO{
		StoragePath: "./image",
	}
	imageService := service.NewImageService(imageDAO, ossImageDAO, locImageDAO)
	imageController := controller.NewImageController(imageService)

	// 注册路由
	r.POST("/api/register", userController.UserRegister)
	r.POST("/api/login", userController.UserLogin)
	r.POST("/api/get", imageController.GetByUsername)
	r.GET("/api/image/:id", imageController.GetByID)

	auth := r.Group("/api")
	auth.Use(userController.AuthMiddleware)
	auth.POST("/upload", imageController.CreateImage)
}
