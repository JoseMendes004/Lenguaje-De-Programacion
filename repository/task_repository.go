package repository

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/usuario/task-manager/models"
)

// TaskRepository maneja la persistencia de tareas
type TaskRepository struct {
	filePath string
	mu       sync.RWMutex
}

// NewTaskRepository crea una nueva instancia del repositorio
func NewTaskRepository(filePath string) *TaskRepository {
	return &TaskRepository{
		filePath: filePath,
	}
}

// SaveTasks guarda todas las tareas en el archivo JSON
func (r *TaskRepository) SaveTasks(tasks map[string]*models.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Convertir el map a slice para mejor formato JSON
	taskList := make([]*models.Task, 0, len(tasks))
	for _, task := range tasks {
		taskList = append(taskList, task)
	}

	data, err := json.MarshalIndent(taskList, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0644)
}

// LoadTasks carga todas las tareas desde el archivo JSON
func (r *TaskRepository) LoadTasks() (map[string]*models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Si el archivo no existe, retornar un map vacío
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return make(map[string]*models.Task), nil
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var taskList []*models.Task
	if err := json.Unmarshal(data, &taskList); err != nil {
		return nil, err
	}

	// Convertir slice a map para acceso rápido por ID
	tasks := make(map[string]*models.Task)
	for _, task := range taskList {
		tasks[task.ID] = task
	}

	return tasks, nil
}
