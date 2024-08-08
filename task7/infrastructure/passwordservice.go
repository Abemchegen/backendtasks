package infrastructure

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(Password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func Compare(password1 string, password2 string) error {

	if err := bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2)); err != nil {
		return err
	}
	return nil
}
