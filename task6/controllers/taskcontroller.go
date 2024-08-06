package controllers

import (
	"net/http"
	"task6/data"
	"task6/models"

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
	role, exists := ctx.Get("role")
	if !exists || role != "user" {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "user-only"})
		return
	}
	var newtask models.Task

	if err := ctx.BindJSON(&newtask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userid, exists := ctx.Get("user_id")

	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user id not found"})
		return
	}
	useridstr, ok := userid.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user ID is not a valid string"})
		return
	}
	userObjectID, err := primitive.ObjectIDFromHex(useridstr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user ID is not a valid ObjectID"})
		return
	}
	newtask.UserID = userObjectID

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
	role, exists := ctx.Get("role")
	if !exists || role != "user" {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "user-only"})
		return
	}
	id := ctx.Param("id")
	task, err := tc.service.GetTask(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, task)

}
func (tc *TaskController) GetTasks(ctx *gin.Context) {
	role, exists := ctx.Get("role")
	if !exists || role != "user" {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "user-only"})
		return
	}
	userid, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "user id doesn't exsist"})
		return
	}
	userID := userid.(string)
	tasks, err := tc.service.GetTasks(userID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tasks)

}
func (tc *TaskController) UpdateTask(ctx *gin.Context) {
	role, exists := ctx.Get("role")
	if !exists || role != "user" {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "user-only"})
		return
	}

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
	role, exists := ctx.Get("role")
	if !exists || role != "user" {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "user-only"})
		return
	}

	id := ctx.Param("id")
	result, err := tc.service.RemoveTask(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": result})

}
