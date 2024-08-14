package usecases

import (
	"task8/domain"
	"task8/infrastructure"
)

type UserUsecase struct {
	repository domain.UserRepositoryInterface
	js         infrastructure.JWTService
}

func NewUserUsecase(repository domain.UserRepositoryInterface, js infrastructure.JWTService) *UserUsecase {
	return &UserUsecase{repository: repository, js: js}
}

func (us *UserUsecase) Register(user *domain.User) error {

	err := us.repository.Register(user)

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
	token, err := us.js.NewToken(user.ID.Hex(), user.Email, role)
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
