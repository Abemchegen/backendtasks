package controllers

import (
	"net/http"
	"task6/data"
	"task6/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	service data.UserService
}

func NewUserController(service data.UserService) *UserController {
	return &UserController{service: service}
}

func (us *UserController) Register(ctx *gin.Context) {

	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := us.service.Register(&user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "user": user})
		return
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)

	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrive the inserted ID"})
		return
	}
	user.ID = oid
	ctx.JSON(http.StatusCreated, gin.H{"message": "registered successfully", "user": user})

}

func (us *UserController) Login(ctx *gin.Context) {

	var loginUser models.User
	if err := ctx.ShouldBindJSON(&loginUser); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := us.service.Login(&loginUser)

	if err != nil || token == "" {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "user logged in successfully", "token": token})

}

func (us *UserController) GetUser(ctx *gin.Context) {

	role, exists := ctx.Get("role")

	if !exists || role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "admin-only", "role": role})
		return
	}

	email := ctx.Param("email")
	user, err := us.service.GetUser(email)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)

}

func (us *UserController) GetUsers(ctx *gin.Context) {
	role, exists := ctx.Get("role")

	if !exists || role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "admin-only"})
	}

	users, err := us.service.GetUsers()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)

}
