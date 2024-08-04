package controllers

import (
	"net/http"
	"task5/data"
	"task5/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskController struct {
	service *data.TaskService
}

func NewTaskController(service *data.TaskService) *TaskController {
	return &TaskController{service: service}
}

func (tc *TaskController) CreateTask(ctx *gin.Context) {
	var newtask models.Task

	if err := ctx.BindJSON(&newtask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := tc.service.CreateTask(&newtask)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		newtask.ID = oid
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrive the inserted ID"})
		return
	}

	ctx.JSON(http.StatusCreated, newtask)

}
func (tc *TaskController) GetTask(ctx *gin.Context) {
	id := ctx.Param("id")
	task, err := tc.service.GetTask(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, task)

}
func (tc *TaskController) GetTasks(ctx *gin.Context) {

	tasks, err := tc.service.GetTasks()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tasks)

}
func (tc *TaskController) UpdateTask(ctx *gin.Context) {

	id := ctx.Param("id")

	var updatedTask models.Task

	if err := ctx.BindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := tc.service.UpdateTask(id, &updatedTask)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, task)

}
func (tc *TaskController) RemoveTask(ctx *gin.Context) {

	id := ctx.Param("id")
	result, err := tc.service.RemoveTask(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": result})

}
