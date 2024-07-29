package routes

import (
	grpclient "todowithgin/grpcClient"
	"todowithgin/handlers"
	// "todowithgin/middleware"
	"todowithgin/service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {

	dbService := service.NewDBservice()
	grpcService := &grpclient.GrpcService{}
	// middleware := &middleware.MiddlewareService{}
	taskHandler := handlers.NewTaskHandler(dbService, grpcService)
	router := gin.Default()

	router.GET("/tasks", taskHandler.GetTasksHandler)
	router.GET("/tasks/:id", taskHandler.GetTaskHandler)

	r1 := router.Group("/tasks", taskHandler.PrivateRouteHandler)
	{
		r1.POST("/", taskHandler.CreateTaskHandler)
		r1.PUT("/:id", taskHandler.UpdateTaskHandler)
		r1.DELETE("/:id", taskHandler.DeleteTaskHandler)

	}

	router.POST("/signup", taskHandler.CreateUserHandler)
	router.POST("/login", taskHandler.LoginHandler)

	return router

	// router.GET("/test", taskHandler.Test)
	// router.GET("/auth", middleware.ProtectedRoutes, func(ctx *gin.Context) {
	// 	user,_:=ctx.Get("username")
	// 	userId,_:=ctx.Get("userId")
	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"username":user,
	// 		"userId":userId,
	// 	})
	// })

}
