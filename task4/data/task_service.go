package data

import (
	"sync"
	"task4/models"
)

var (
	count int
	tasks = map[int]models.Task{}
	mu    sync.Mutex
)

func GetAllTasks() map[int]models.Task {
	mu.Lock()
	defer mu.Unlock()
	return tasks
}

func GetTask(id int) (models.Task, bool) {
	mu.Lock()
	defer mu.Unlock()
	task, exists := tasks[id]
	return task, exists
}
func UpdateTask(id int, updatedTask models.Task) {
	mu.Lock()
	defer mu.Unlock()
	tasks[id] = updatedTask
}

func DeleteTask(id int) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := tasks[id]; exists {
		delete(tasks, id)
		return true
	}
	return false
}

func CreateTask(newTask models.Task) int {
	mu.Lock()
	defer mu.Unlock()
	count++
	newTask.ID = count
	tasks[count] = newTask
	return count
}
