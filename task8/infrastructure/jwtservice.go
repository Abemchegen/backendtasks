package infrastructure

import (
	"os"

	"github.com/dgrijalva/jwt-go"
)

type JWTService interface {
	NewToken(id string, email string, role string) (string, error)
}

func NewToken(id string, email string, role string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
	})

	var jwtSecret []byte = []byte(os.Getenv("secret"))

	jwtToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}
