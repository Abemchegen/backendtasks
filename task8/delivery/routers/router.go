package routers

import (
	"task7/delivery/controllers"
	"task7/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetRouter(c *controllers.TaskController, u *controllers.UserController) *gin.Engine {

	router := gin.Default()
	route := router.Group("/", infrastructure.UserAuthorizaiton())
	{

		route.GET("tasks/", c.GetTasks)
		route.GET("tasks/:id", c.GetTask)
		route.PUT("tasks/:id", c.UpdateTask)
		route.POST("tasks/", c.CreateTask)
		route.DELETE("tasks/:id", c.RemoveTask)
		route.GET("users/", u.GetUsers)
		route.GET("user/:email", u.GetUser)
	}
	router.POST("/register", u.Register)
	router.POST("/login", u.Login)

	return router

}
