package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seltGrowth/internal/growth_record/controller"
	"seltGrowth/internal/growth_record/middleware"
)

func Router() http.Handler {
	router := gin.Default()
	helloHandler := controller.NewHelloHandler()
	activityController := controller.NewActivityController()
	taskController := controller.NewTaskController()
	userController := controller.NewUserController()
	labelController := controller.NewLabelController();

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
			task.POST("/complete", taskController.Complete)
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

	return router
}
