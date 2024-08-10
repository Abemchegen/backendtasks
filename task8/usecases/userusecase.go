package usecases

import (
	"errors"
	"task7/domain"
	"task7/infrastructure"
)

type UserUsecase struct {
	repository domain.UserRepositoryInterface
}

func NewUserUsecase(repository domain.UserRepositoryInterface) *UserUsecase {
	return &UserUsecase{repository: repository}
}

func (us *UserUsecase) Register(user *domain.User) error {

	if user.Email == "" || user.Role == "" {
		return errors.New("incomplete information")
	}
	hashedPassword, err := infrastructure.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	err = us.repository.Register(user)

	if err != nil {
		return err
	}
	return nil
}

func (us *UserUsecase) Login(user *domain.User) (string, error) {

	role, err := us.repository.Login(user)
	if err != nil {
		return "", err
	}
	token, err := infrastructure.NewToken(user.ID.Hex(), user.Email, role)
	if err != nil {
		return "", err
	}
	return token, nil

}

func (us *UserUsecase) GetUser(email string) (*domain.User, error) {

	user, err := us.repository.GetUser(email)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (us *UserUsecase) GetUsers() (*[]domain.User, error) {

	users, err := us.repository.GetUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}
