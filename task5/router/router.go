package router

import (
	"task5/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(c *controllers.TaskController) *gin.Engine {
	router := gin.Default()

	router.GET("tasks/", c.GetTasks)
	router.GET("tasks/:id", c.GetTask)
	router.PUT("tasks/:id", c.UpdateTask)
	router.POST("tasks/", c.CreateTask)
	router.DELETE("tasks/:id", c.RemoveTask)

	return router

}
