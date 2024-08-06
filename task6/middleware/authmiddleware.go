package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func UserAuthorizaiton() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		authheader := ctx.GetHeader("Authorization")
		if authheader == "" {
			ctx.JSON(401, gin.H{"error": "Authorization header is required"})
			ctx.Abort()
			return
		}
		authstring := strings.Split(authheader, " ")

		if len(authstring) != 2 || strings.ToLower(authstring[0]) != "bearer" {
			ctx.JSON(401, gin.H{"error": "Invalid authorization header"})
			ctx.Abort()
			return
		}

		secret := []byte(os.Getenv("secret"))

		token, err := jwt.Parse(authstring[1], func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil

		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Failed to parse token: %v", err)})
			ctx.Abort()
			return
		}
		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT"})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}
		userID, userIDExists := claims["user_id"]
		email, emailExists := claims["email"]
		role, roleExists := claims["role"]

		if !userIDExists || !emailExists || !roleExists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token does not contain necessary claims"})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", userID)
		ctx.Set("email", email)
		ctx.Set("role", role)
		ctx.Next()
	}

}
