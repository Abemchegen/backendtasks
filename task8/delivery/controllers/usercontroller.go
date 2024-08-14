package controllers

import (
	"errors"
	"net/http"
	"task8/domain"

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
	if user.Email == "" || user.Role == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New("incomplete information").Error()})
		return
	}

	err := us.usecase.Register(&user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "registered successfully"})

}

func (us *UserController) Login(ctx *gin.Context) {
	var loginUser domain.User

	if err := ctx.ShouldBindJSON(&loginUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if loginUser.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "incomplete information"})
		return
	}

	token, err := us.usecase.Login(&loginUser)

	if err != nil || token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "token": token})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user logged in successfully", "token": token})
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
