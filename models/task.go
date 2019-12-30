package models

import (
	"encoding/json"
	"time"
)

const (
	// TaskCopy const.
	TaskCopy = 0
	// TaskDelete const.
	TaskDelete = 1
	// TaskMove const.
	TaskMove = 2

	// StatusProgress const.
	StatusProgress = 0
	// StatusEnded const.
	StatusEnded = 1
)

// Task struct definition.
type Task struct {
	Type      uint8
	Name      string
	StartedAt time.Time
	EndedAt   time.Time
	Status    uint8
	Error     error
}

// End task.
func (task *Task) End(err error) error {
	task.Status = StatusEnded
	task.EndedAt = time.Now()
	task.Error = err
	return err
}

// NewTask constructor.
func NewTask(name string, taskType uint8) Task {
	task := Task{Name: name, StartedAt: time.Now(), Status: StatusProgress}
	return task
}

// NewCopyTask create new copy task.
func NewCopyTask(src, dst string) Task {
	return NewTask("Copy "+src+" to "+dst, TaskCopy)
}

// NewMoveTask create new move task.
func NewMoveTask(src, dst string) Task {
	return NewTask("Move "+src+" to "+dst, TaskMove)
}

// MarshalJSON function.
func (task *Task) MarshalJSON() (text []byte, err error) {
	s := struct {
		Name      string    `json:"name"`
		Status    string    `json:"status"`
		StartedAt time.Time `json:"startedAt"`
		EndedAt   time.Time `json:"endedAt"`
		Error     error     `json:"error"`
	}{
		Name:      task.Name,
		StartedAt: task.StartedAt,
		EndedAt:   task.EndedAt,
		Error:     task.Error,
	}

	switch task.Status {
	case StatusProgress:
		s.Status = "progress"
	case StatusEnded:
		s.Status = "ended"
	}

	return json.Marshal(s)
}
