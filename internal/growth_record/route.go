package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seltGrowth/internal/growth_record/controller"
	"seltGrowth/internal/growth_record/middleware"
)

func router() http.Handler {
	router := gin.Default()
	helloHandler := controller.NewHelloHandler()
	phoneUseController := controller.NewPhoneUseController()
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

		phoneUser := v1.Group("/phone")
		{
			phoneUser.POST("/useRecord", phoneUseController.UploadRecord)
			phoneUser.GET("/overview", phoneUseController.Overview)
			phoneUser.GET("/activityHistory", phoneUseController.ActivityHistory)
		}

		activity := v1.Group("/activity")
		{
			activity.GET("/list", activityController.GetActivities)
		}

		task := v1.Group("/task")
		{
			task.GET("/list", taskController.TaskList)
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
