package routers

import (
	"task6/controllers"
	"task6/middleware"

	"github.com/gin-gonic/gin"
)

func SetRouter(c *controllers.TaskController, u *controllers.UserController) *gin.Engine {

	router := gin.Default()

	router.GET("tasks/", middleware.UserAuthorizaiton(), c.GetTasks)
	router.GET("tasks/:id", middleware.UserAuthorizaiton(), c.GetTask)
	router.PUT("tasks/:id", middleware.UserAuthorizaiton(), c.UpdateTask)
	router.POST("tasks/", middleware.UserAuthorizaiton(), c.CreateTask)
	router.DELETE("tasks/:id", middleware.UserAuthorizaiton(), c.RemoveTask)

	router.GET("users/", middleware.UserAuthorizaiton(), u.GetUsers)
	router.GET("user/:email", middleware.UserAuthorizaiton(), u.GetUser)
	router.POST("/register", u.Register)
	router.POST("/login", u.Login)

	return router

}
