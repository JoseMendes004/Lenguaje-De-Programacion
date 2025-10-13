package models

import (
	"errors"
	"time"
)

// TaskStatus representa el estado de una tarea usando un enum
type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusCancelled  TaskStatus = "cancelled"
)

// Validar si el estado es válido
func (s TaskStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusInProgress, StatusCompleted, StatusCancelled:
		return true
	}
	return false
}

// Task representa una tarea en el sistema
type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Errores personalizados
var (
	ErrEmptyTitle       = errors.New("el título no puede estar vacío")
	ErrInvalidStatus    = errors.New("estado de tarea inválido")
	ErrTaskNotFound     = errors.New("tarea no encontrada")
	ErrInvalidID        = errors.New("ID de tarea inválido")
)

// NewTask crea una nueva tarea con validaciones
func NewTask(id, title, description string) (*Task, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}
	if id == "" {
		return nil, ErrInvalidID
	}

	now := time.Now()
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// UpdateStatus actualiza el estado de la tarea
func (t *Task) UpdateStatus(status TaskStatus) error {
	if !status.IsValid() {
		return ErrInvalidStatus
	}
	t.Status = status
	t.UpdatedAt = time.Now()
	return nil
}

// Update actualiza los campos de la tarea
func (t *Task) Update(title, description string) error {
	if title == "" {
		return ErrEmptyTitle
	}
	t.Title = title
	t.Description = description
	t.UpdatedAt = time.Now()
	return nil
}
