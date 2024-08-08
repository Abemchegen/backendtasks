package controllers

import (
	"net/http"
	"task7/domain"
	"task7/usecases"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	usecase *usecases.TaskUsecase
}

func NewTaskController(usecase *usecases.TaskUsecase) *TaskController {
	return &TaskController{usecase: usecase}
}

func (tc *TaskController) CreateTask(ctx *gin.Context) {

	role, exists := ctx.Get("role")
	if !exists || role != "user" {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "user-only"})
		return
	}
	var newtask domain.Task

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
	err := tc.usecase.CreateTask(&newtask, useridstr)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	task, err := tc.usecase.GetTask(id)

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
	tasks, err := tc.usecase.GetTasks(userID)

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

	var updatedTask domain.Task

	if err := ctx.BindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := tc.usecase.UpdateTask(id, &updatedTask)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedTask)

}
func (tc *TaskController) RemoveTask(ctx *gin.Context) {
	role, exists := ctx.Get("role")
	if !exists || role != "user" {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "user-only"})
		return
	}

	id := ctx.Param("id")
	err := tc.usecase.RemoveTask(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "task removed"})

}
