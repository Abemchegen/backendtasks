package usecases

import (
	"errors"
	"task8/domain"
)

type TaskUsecase struct {
	repository domain.TaskRepositoryInterface
}

func NewTaskUsecase(repository domain.TaskRepositoryInterface) *TaskUsecase {
	return &TaskUsecase{repository: repository}
}

func (tc *TaskUsecase) CreateTask(newtask *domain.Task, userid string) error {

	if newtask.Description == "" || newtask.Status == "" || newtask.Title == "" {
		return errors.New("incomplete information")
	}
	return tc.repository.CreateTask(newtask, userid)
}

func (tc *TaskUsecase) GetTask(id string) (*domain.Task, error) {
	return tc.repository.GetTask(id)
}

func (tc *TaskUsecase) GetTasks(userID string) (*[]domain.Task, error) {
	return tc.repository.GetTasks(userID)
}

func (tc *TaskUsecase) UpdateTask(id string, updatedTask *domain.Task) error {
	return tc.repository.UpdateTask(id, updatedTask)
}

func (tc *TaskUsecase) RemoveTask(id string) error {
	return tc.repository.RemoveTask(id)
}
