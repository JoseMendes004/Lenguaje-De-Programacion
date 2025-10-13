package tests

import (
	"testing"

	"github.com/usuario/task-manager/models"
)

func TestNewTask(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		title       string
		description string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "tarea válida",
			id:          "1",
			title:       "Tarea de prueba",
			description: "Descripción de prueba",
			wantErr:     false,
		},
		{
			name:        "título vacío",
			id:          "2",
			title:       "",
			description: "Descripción",
			wantErr:     true,
			expectedErr: models.ErrEmptyTitle,
		},
		{
			name:        "ID vacío",
			id:          "",
			title:       "Título",
			description: "Descripción",
			wantErr:     true,
			expectedErr: models.ErrInvalidID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := models.NewTask(tt.id, tt.title, tt.description)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Se esperaba un error pero no se obtuvo ninguno")
				}
				if err != tt.expectedErr {
					t.Errorf("Error esperado: %v, obtenido: %v", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("No se esperaba error pero se obtuvo: %v", err)
				return
			}

			if task.ID != tt.id {
				t.Errorf("ID esperado: %s, obtenido: %s", tt.id, task.ID)
			}
			if task.Title != tt.title {
				t.Errorf("Título esperado: %s, obtenido: %s", tt.title, task.Title)
			}
			if task.Status != models.StatusPending {
				t.Errorf("Estado esperado: %s, obtenido: %s", models.StatusPending, task.Status)
			}
		})
	}
}

func TestTaskUpdateStatus(t *testing.T) {
	task, _ := models.NewTask("1", "Tarea", "Descripción")

	tests := []struct {
		name        string
		status      models.TaskStatus
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "estado válido - in_progress",
			status:  models.StatusInProgress,
			wantErr: false,
		},
		{
			name:    "estado válido - completed",
			status:  models.StatusCompleted,
			wantErr: false,
		},
		{
			name:        "estado inválido",
			status:      models.TaskStatus("invalid"),
			wantErr:     true,
			expectedErr: models.ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := task.UpdateStatus(tt.status)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Se esperaba un error pero no se obtuvo ninguno")
				}
				if err != tt.expectedErr {
					t.Errorf("Error esperado: %v, obtenido: %v", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("No se esperaba error pero se obtuvo: %v", err)
			}

			if task.Status != tt.status {
				t.Errorf("Estado esperado: %s, obtenido: %s", tt.status, task.Status)
			}
		})
	}
}

func TestTaskUpdate(t *testing.T) {
	task, _ := models.NewTask("1", "Título original", "Descripción original")

	tests := []struct {
		name        string
		title       string
		description string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "actualización válida",
			title:       "Nuevo título",
			description: "Nueva descripción",
			wantErr:     false,
		},
		{
			name:        "título vacío",
			title:       "",
			description: "Descripción",
			wantErr:     true,
			expectedErr: models.ErrEmptyTitle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := task.Update(tt.title, tt.description)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Se esperaba un error pero no se obtuvo ninguno")
				}
				if err != tt.expectedErr {
					t.Errorf("Error esperado: %v, obtenido: %v", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("No se esperaba error pero se obtuvo: %v", err)
			}

			if task.Title != tt.title {
				t.Errorf("Título esperado: %s, obtenido: %s", tt.title, task.Title)
			}
			if task.Description != tt.description {
				t.Errorf("Descripción esperada: %s, obtenida: %s", tt.description, task.Description)
			}
		})
	}
}

func TestTaskStatusIsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   models.TaskStatus
		expected bool
	}{
		{"pending válido", models.StatusPending, true},
		{"in_progress válido", models.StatusInProgress, true},
		{"completed válido", models.StatusCompleted, true},
		{"cancelled válido", models.StatusCancelled, true},
		{"estado inválido", models.TaskStatus("invalid"), false},
		{"estado vacío", models.TaskStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.IsValid()
			if result != tt.expected {
				t.Errorf("Resultado esperado: %v, obtenido: %v", tt.expected, result)
			}
		})
	}
}
