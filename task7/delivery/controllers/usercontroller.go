package controllers

import (
	"net/http"
	"task7/domain"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	usecase domain.UserUsecaseInterface
}

func NewUserController(usecase domain.UserUsecaseInterface) *UserController {
	return &UserController{usecase: usecase}
}

func (us *UserController) Register(ctx *gin.Context) {

	var user domain.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := us.usecase.Register(&user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "registered successfully", "user": user})

}

func (us *UserController) Login(ctx *gin.Context) {

	var loginUser domain.User
	if err := ctx.ShouldBindJSON(&loginUser); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := us.usecase.Login(&loginUser)

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
	user, err := us.usecase.GetUser(email)

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

	users, err := us.usecase.GetUsers()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)

}
