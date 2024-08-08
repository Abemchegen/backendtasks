package infrastructure

import "golang.org/x/crypto/bcrypt"

func Hash(Password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
