package service

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/usuario/task-manager/models"
	"github.com/usuario/task-manager/repository"
)

// TaskService maneja la lógica de negocio de las tareas
type TaskService struct {
	repo       *repository.TaskRepository
	tasks      map[string]*models.Task
	nextID     int
	mu         sync.RWMutex
}

// NewTaskService crea una nueva instancia del servicio
func NewTaskService(repo *repository.TaskRepository) (*TaskService, error) {
	service := &TaskService{
		repo:   repo,
		tasks:  make(map[string]*models.Task),
		nextID: 1,
	}

	// Cargar tareas existentes
	if err := service.loadTasks(); err != nil {
		return nil, fmt.Errorf("error al cargar tareas: %w", err)
	}

	return service, nil
}

// loadTasks carga las tareas desde el repositorio
func (s *TaskService) loadTasks() error {
	tasks, err := s.repo.LoadTasks()
	if err != nil {
		return err
	}

	s.tasks = tasks

	// Calcular el siguiente ID
	maxID := 0
	for id := range s.tasks {
		if numID, err := strconv.Atoi(id); err == nil && numID > maxID {
			maxID = numID
		}
	}
	s.nextID = maxID + 1

	return nil
}

// saveTasks guarda las tareas en el repositorio
func (s *TaskService) saveTasks() error {
	return s.repo.SaveTasks(s.tasks)
}

// CreateTask crea una nueva tarea
func (s *TaskService) CreateTask(title, description string) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := strconv.Itoa(s.nextID)
	task, err := models.NewTask(id, title, description)
	if err != nil {
		return nil, err
	}

	s.tasks[id] = task
	s.nextID++

	if err := s.saveTasks(); err != nil {
		return nil, fmt.Errorf("error al guardar tarea: %w", err)
	}

	return task, nil
}

// GetTask obtiene una tarea por ID
func (s *TaskService) GetTask(id string) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, models.ErrTaskNotFound
	}

	return task, nil
}

// GetAllTasks obtiene todas las tareas
func (s *TaskService) GetAllTasks() []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// UpdateTask actualiza una tarea existente
func (s *TaskService) UpdateTask(id, title, description string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return models.ErrTaskNotFound
	}

	if err := task.Update(title, description); err != nil {
		return err
	}

	if err := s.saveTasks(); err != nil {
		return fmt.Errorf("error al guardar cambios: %w", err)
	}

	return nil
}

// UpdateTaskStatus actualiza el estado de una tarea
func (s *TaskService) UpdateTaskStatus(id string, status models.TaskStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return models.ErrTaskNotFound
	}

	if err := task.UpdateStatus(status); err != nil {
		return err
	}

	if err := s.saveTasks(); err != nil {
		return fmt.Errorf("error al guardar cambios: %w", err)
	}

	return nil
}

// DeleteTask elimina una tarea
func (s *TaskService) DeleteTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return models.ErrTaskNotFound
	}

	delete(s.tasks, id)

	if err := s.saveTasks(); err != nil {
		return fmt.Errorf("error al guardar cambios: %w", err)
	}

	return nil
}

// GetTasksByStatus obtiene tareas filtradas por estado
func (s *TaskService) GetTasksByStatus(status models.TaskStatus) []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*models.Task, 0)
	for _, task := range s.tasks {
		if task.Status == status {
			tasks = append(tasks, task)
		}
	}

	return tasks
}
