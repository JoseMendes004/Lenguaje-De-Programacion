package tests

import (
	"os"
	"testing"

	"github.com/usuario/task-manager/models"
	"github.com/usuario/task-manager/repository"
	"github.com/usuario/task-manager/service"
)

func setupTestService(t *testing.T) (*service.TaskService, string) {
	// Crear archivo temporal para pruebas
	tmpFile := "test_tasks.json"
	
	// Limpiar archivo si existe
	os.Remove(tmpFile)
	
	repo := repository.NewTaskRepository(tmpFile)
	svc, err := service.NewTaskService(repo)
	if err != nil {
		t.Fatalf("Error al crear servicio: %v", err)
	}
	
	return svc, tmpFile
}

func cleanupTestService(tmpFile string) {
	os.Remove(tmpFile)
}

func TestCreateTask(t *testing.T) {
	svc, tmpFile := setupTestService(t)
	defer cleanupTestService(tmpFile)

	task, err := svc.CreateTask("Tarea de prueba", "Descripción de prueba")
	if err != nil {
		t.Fatalf("Error al crear tarea: %v", err)
	}

	if task.Title != "Tarea de prueba" {
		t.Errorf("Título esperado: 'Tarea de prueba', obtenido: '%s'", task.Title)
	}

	if task.Status != models.StatusPending {
		t.Errorf("Estado esperado: %s, obtenido: %s", models.StatusPending, task.Status)
	}
}

func TestGetTask(t *testing.T) {
	svc, tmpFile := setupTestService(t)
	defer cleanupTestService(tmpFile)

	// Crear tarea
	created, _ := svc.CreateTask("Tarea", "Descripción")

	// Obtener tarea
	task, err := svc.GetTask(created.ID)
	if err != nil {
		t.Fatalf("Error al obtener tarea: %v", err)
	}

	if task.ID != created.ID {
		t.Errorf("ID esperado: %s, obtenido: %s", created.ID, task.ID)
	}

	// Intentar obtener tarea inexistente
	_, err = svc.GetTask("999")
	if err != models.ErrTaskNotFound {
		t.Errorf("Error esperado: %v, obtenido: %v", models.ErrTaskNotFound, err)
	}
}

func TestUpdateTask(t *testing.T) {
	svc, tmpFile := setupTestService(t)
	defer cleanupTestService(tmpFile)

	// Crear tarea
	task, _ := svc.CreateTask("Título original", "Descripción original")

	// Actualizar tarea
	err := svc.UpdateTask(task.ID, "Título actualizado", "Descripción actualizada")
	if err != nil {
		t.Fatalf("Error al actualizar tarea: %v", err)
	}

	// Verificar actualización
	updated, _ := svc.GetTask(task.ID)
	if updated.Title != "Título actualizado" {
		t.Errorf("Título esperado: 'Título actualizado', obtenido: '%s'", updated.Title)
	}
}

func TestUpdateTaskStatus(t *testing.T) {
	svc, tmpFile := setupTestService(t)
	defer cleanupTestService(tmpFile)

	// Crear tarea
	task, _ := svc.CreateTask("Tarea", "Descripción")

	// Actualizar estado
	err := svc.UpdateTaskStatus(task.ID, models.StatusCompleted)
	if err != nil {
		t.Fatalf("Error al actualizar estado: %v", err)
	}

	// Verificar actualización
	updated, _ := svc.GetTask(task.ID)
	if updated.Status != models.StatusCompleted {
		t.Errorf("Estado esperado: %s, obtenido: %s", models.StatusCompleted, updated.Status)
	}
}

func TestDeleteTask(t *testing.T) {
	svc, tmpFile := setupTestService(t)
	defer cleanupTestService(tmpFile)

	// Crear tarea
	task, _ := svc.CreateTask("Tarea", "Descripción")

	// Eliminar tarea
	err := svc.DeleteTask(task.ID)
	if err != nil {
		t.Fatalf("Error al eliminar tarea: %v", err)
	}

	// Verificar que no existe
	_, err = svc.GetTask(task.ID)
	if err != models.ErrTaskNotFound {
		t.Errorf("Error esperado: %v, obtenido: %v", models.ErrTaskNotFound, err)
	}
}

func TestGetTasksByStatus(t *testing.T) {
	svc, tmpFile := setupTestService(t)
	defer cleanupTestService(tmpFile)

	// Crear varias tareas
	svc.CreateTask("Tarea 1", "Descripción 1")
	task2, _ := svc.CreateTask("Tarea 2", "Descripción 2")
	task3, _ := svc.CreateTask("Tarea 3", "Descripción 3")

	// Cambiar estados
	svc.UpdateTaskStatus(task2.ID, models.StatusCompleted)
	svc.UpdateTaskStatus(task3.ID, models.StatusCompleted)

	// Obtener tareas completadas
	completed := svc.GetTasksByStatus(models.StatusCompleted)
	if len(completed) != 2 {
		t.Errorf("Se esperaban 2 tareas completadas, se obtuvieron: %d", len(completed))
	}

	// Obtener tareas pendientes
	pending := svc.GetTasksByStatus(models.StatusPending)
	if len(pending) != 1 {
		t.Errorf("Se esperaba 1 tarea pendiente, se obtuvieron: %d", len(pending))
	}
}

func TestPersistence(t *testing.T) {
	tmpFile := "test_persistence.json"
	defer cleanupTestService(tmpFile)

	// Crear servicio y agregar tareas
	{
		repo := repository.NewTaskRepository(tmpFile)
		svc, _ := service.NewTaskService(repo)
		svc.CreateTask("Tarea persistente", "Esta tarea debe persistir")
	}

	// Crear nuevo servicio y verificar que las tareas se cargaron
	{
		repo := repository.NewTaskRepository(tmpFile)
		svc, _ := service.NewTaskService(repo)
		tasks := svc.GetAllTasks()
		
		if len(tasks) != 1 {
			t.Errorf("Se esperaba 1 tarea persistida, se obtuvieron: %d", len(tasks))
		}
		
		if len(tasks) > 0 && tasks[0].Title != "Tarea persistente" {
			t.Errorf("Título esperado: 'Tarea persistente', obtenido: '%s'", tasks[0].Title)
		}
	}
}
