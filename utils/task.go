package utils

import (
	"filebrowser/models"

	"github.com/google/uuid"
)

// Tasks utils.
type Tasks struct {
	list map[uuid.UUID]*models.Task
}

// NewTasks constructor.
func NewTasks() Tasks {
	tasks := Tasks{}
	tasks.list = make(map[uuid.UUID]*models.Task)
	return tasks
}

// Add task.
func (tasks *Tasks) Add(task *models.Task) uuid.UUID {
	id := uuid.New()
	tasks.list[id] = task
	return id
}

// Get single task.
func (tasks *Tasks) Get(id uuid.UUID) (*models.Task, bool) {
	task, finded := tasks.list[id]
	return task, finded
}

// All tasks.
func (tasks *Tasks) All() []*models.Task {
	values := []*models.Task{}
	for _, value := range tasks.list {
		values = append(values, value)
	}
	return values
}

// Remove task.
func (tasks *Tasks) Remove(id uuid.UUID) {
	delete(tasks.list, id)
}

// CleanEnded all ended tasks.
func (tasks *Tasks) CleanEnded() {
	for k, v := range tasks.list {
		if v.Status == models.StatusEnded {
			tasks.Remove(k)
		}
	}
}

// Clean all tasks.
func (tasks *Tasks) Clean() {
	for k := range tasks.list {
		tasks.Remove(k)
	}
}
