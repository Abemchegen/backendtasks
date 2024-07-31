package controllers

import (
	"net/http"
	"strconv"
	"time"

	"task4/data"
	"task4/models"

	"github.com/gin-gonic/gin"
)

func GetTasks(ctx *gin.Context) {
	tasks := data.GetAllTasks()
	ctx.JSON(http.StatusOK, tasks)
}

func GetTaskByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	task, exists := data.GetTask(id)
	if exists {
		ctx.JSON(http.StatusOK, task)
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
	}
}

func UpdateTask(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var updatedTask models.Task
	if err := ctx.BindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, exists := data.GetTask(id)
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if updatedTask.Title != "" {
		task.Title = updatedTask.Title
	}
	if updatedTask.Description != "" {
		task.Description = updatedTask.Description
	}
	if updatedTask.Status != "" {
		task.Status = updatedTask.Status
	}

	data.UpdateTask(id, task)
	ctx.JSON(http.StatusOK, task)
}

func DeleteTask(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if !data.DeleteTask(id) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task removed"})
}

func CreateTask(ctx *gin.Context) {
	var newTask models.Task

	if err := ctx.BindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tid := data.CreateTask(newTask)
	newTask.DueDate = time.Now()
	newTask.ID = tid
	ctx.JSON(http.StatusCreated, gin.H{"task": newTask})
}
