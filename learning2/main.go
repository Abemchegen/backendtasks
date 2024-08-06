package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var users = make(map[string]*User)
var jwtSecret = []byte("your_jwt_secret")

func main() {

	router := gin.Default()
	router.POST("/register", func(ctx *gin.Context) {
		var user User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		user.Password = string(hashedPassword)
		users[user.Email] = &user

		ctx.JSON(200, gin.H{"message": "registered successfully"})
	})

	router.POST("/login", func(ctx *gin.Context) {
		var loginUser User
		if err := ctx.ShouldBindJSON(&loginUser); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Attempting login for user: %s", loginUser.Email)
		u, exists := users[loginUser.Email]
		if !exists {
			log.Printf("User not found: %s", loginUser.Email)
			ctx.JSON(401, gin.H{"error": "invalid email or password"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(loginUser.Password)); err != nil {
			log.Printf("Password mismatch for user: %s", loginUser.Email)
			ctx.JSON(401, gin.H{"error": "invalid email or password"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": u.ID,
			"email":   u.Email,
		})
		jwtToken, err := token.SignedString(jwtSecret)
		if err != nil {
			log.Printf("Error signing token for user: %s, error: %v", loginUser.Email, err)
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		log.Printf("User logged in successfully: %s", loginUser.Email)
		ctx.JSON(200, gin.H{"message": "user logged in successfully", "token": jwtToken})
	})

	router.GET("/secure", AuthMiddleware(), func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "This is a secure route"})
	})

	router.GET("/users", AuthMiddleware(), func(ctx *gin.Context) {
		userList := make([]User, 0, len(users))
		for _, user := range users {
			userList = append(userList, *user)
		}
		ctx.JSON(200, userList)
	})

	router.Run()
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(401, gin.H{"error": "Authorization header is required"})
			ctx.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			ctx.JSON(401, gin.H{"error": "Invalid authorization header"})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(authParts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(401, gin.H{"error": "Invalid JWT"})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims["user_id"])
		ctx.Set("email", claims["email"])
		ctx.Next()
	}
}
