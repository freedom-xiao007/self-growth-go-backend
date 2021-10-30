package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"seltGrowth/internal/growth_record/controller"
	"seltGrowth/internal/growth_record/middleware"
)

func initMongodb() {
	// Setup the mgm default config
	username := os.Getenv("mongo_user")
	passowrd := os.Getenv("mongo_password")
	host := os.Getenv("mongo_host")
	port := os.Getenv("mongo_port")
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, passowrd, host, port)
	log.Info("mongoURI", mongoURI)
	err := mgm.SetDefaultConfig(nil, "phone_record", options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initMongodb()
	router := gin.Default()
	router.Use(middleware.Cors())
	InitRoute(router)
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	}
}

func InitRoute(router *gin.Engine) {
	helloHandler := controller.NewHelloHandler()
	activityController := controller.NewActivityController()
	taskController := controller.NewTaskController()
	userController := controller.NewUserController()
	labelController := controller.NewLabelController()

	// 路由分组、中间件、认证
	v1 := router.Group("/v1", middleware.JWTAuth())
	{
		hello := v1.Group("/hello")
		{
			hello.GET("", helloHandler.Hello)
		}

		activity := v1.Group("/activity")
		{
			activity.GET("/list", activityController.GetActivities)
			activity.POST("/useRecord", activityController.UploadRecord)
			activity.GET("/overview", activityController.Overview)
			activity.GET("/activityHistory", activityController.ActivityHistory)
			activity.POST("/updateActivityModel", activityController.UpdateActivityModel)
		}

		task := v1.Group("/task")
		{
			task.GET("/list", taskController.TaskList)
			task.POST("/add", taskController.AddTask)
			task.POST("/complete/:id", taskController.Complete)
			task.GET("/history", taskController.History)
			task.POST("/addTaskGroup", taskController.AddTaskGroup)
			task.GET("/listByGroup", taskController.TaskListByGroup)
			task.GET("/overview", taskController.Overview)
			task.POST("/deleteGroup/:name", taskController.DeleteGroup)
		}

		label := v1.Group("/label")
		{
			label.POST("/add", labelController.Add)
			label.GET("/list", labelController.List)
		}
	}

	login := router.Group("/auth")
	{
		user := login.Group("/user")
		{
			user.POST("/login", userController.Login)
		}
	}
}