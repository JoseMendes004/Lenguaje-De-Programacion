package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/usuario/task-manager/models"
	"github.com/usuario/task-manager/service"
)

// CLI maneja la interfaz de línea de comandos
type CLI struct {
	service *service.TaskService
	scanner *bufio.Scanner
}

// NewCLI crea una nueva instancia de la CLI
func NewCLI(service *service.TaskService) *CLI {
	return &CLI{
		service: service,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// Run inicia el loop principal de la CLI
func (c *CLI) Run() {
	fmt.Println("=== Gestor de Tareas ===")
	fmt.Println("Escribe 'help' para ver los comandos disponibles")

	for {
		fmt.Print("\n> ")
		if !c.scanner.Scan() {
			break
		}

		input := strings.TrimSpace(c.scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]

		switch command {
		case "help":
			c.showHelp()
		case "create":
			c.createTask()
		case "list":
			c.listTasks(parts)
		case "view":
			c.viewTask(parts)
		case "update":
			c.updateTask(parts)
		case "status":
			c.updateStatus(parts)
		case "delete":
			c.deleteTask(parts)
		case "exit", "quit":
			fmt.Println("¡Hasta luego!")
			return
		default:
			fmt.Println("Comando no reconocido. Escribe 'help' para ver los comandos disponibles.")
		}
	}
}

func (c *CLI) showHelp() {
	help := `
Comandos disponibles:
  create              - Crear una nueva tarea
  list [status]       - Listar todas las tareas o filtrar por estado
  view <id>           - Ver detalles de una tarea
  update <id>         - Actualizar una tarea
  status <id> <estado> - Cambiar el estado de una tarea
  delete <id>         - Eliminar una tarea
  help                - Mostrar esta ayuda
  exit/quit           - Salir del programa

Estados disponibles:
  pending, in_progress, completed, cancelled
`
	fmt.Println(help)
}

func (c *CLI) createTask() {
	fmt.Print("Título: ")
	c.scanner.Scan()
	title := strings.TrimSpace(c.scanner.Text())

	fmt.Print("Descripción: ")
	c.scanner.Scan()
	description := strings.TrimSpace(c.scanner.Text())

	task, err := c.service.CreateTask(title, description)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("✓ Tarea creada con ID: %s\n", task.ID)
}

func (c *CLI) listTasks(parts []string) {
	var tasks []*models.Task

	if len(parts) > 1 {
		status := models.TaskStatus(parts[1])
		if !status.IsValid() {
			fmt.Println("Error: Estado inválido")
			return
		}
		tasks = c.service.GetTasksByStatus(status)
		fmt.Printf("\n=== Tareas con estado: %s ===\n", status)
	} else {
		tasks = c.service.GetAllTasks()
		fmt.Println("\n=== Todas las Tareas ===")
	}

	if len(tasks) == 0 {
		fmt.Println("No hay tareas para mostrar.")
		return
	}

	for _, task := range tasks {
		fmt.Printf("\n[%s] %s\n", task.ID, task.Title)
		fmt.Printf("  Estado: %s\n", task.Status)
		fmt.Printf("  Creada: %s\n", task.CreatedAt.Format("2006-01-02 15:04"))
	}
}

func (c *CLI) viewTask(parts []string) {
	if len(parts) < 2 {
		fmt.Println("Error: Debes proporcionar el ID de la tarea")
		return
	}

	task, err := c.service.GetTask(parts[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("\n=== Tarea #%s ===\n", task.ID)
	fmt.Printf("Título: %s\n", task.Title)
	fmt.Printf("Descripción: %s\n", task.Description)
	fmt.Printf("Estado: %s\n", task.Status)
	fmt.Printf("Creada: %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Actualizada: %s\n", task.UpdatedAt.Format("2006-01-02 15:04:05"))
}

func (c *CLI) updateTask(parts []string) {
	if len(parts) < 2 {
		fmt.Println("Error: Debes proporcionar el ID de la tarea")
		return
	}

	id := parts[1]

	// Verificar que la tarea existe
	_, err := c.service.GetTask(id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Print("Nuevo título: ")
	c.scanner.Scan()
	title := strings.TrimSpace(c.scanner.Text())

	fmt.Print("Nueva descripción: ")
	c.scanner.Scan()
	description := strings.TrimSpace(c.scanner.Text())

	if err := c.service.UpdateTask(id, title, description); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("✓ Tarea actualizada correctamente")
}

func (c *CLI) updateStatus(parts []string) {
	if len(parts) < 3 {
		fmt.Println("Error: Uso: status <id> <estado>")
		return
	}

	id := parts[1]
	status := models.TaskStatus(parts[2])

	if err := c.service.UpdateTaskStatus(id, status); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("✓ Estado actualizado correctamente")
}

func (c *CLI) deleteTask(parts []string) {
	if len(parts) < 2 {
		fmt.Println("Error: Debes proporcionar el ID de la tarea")
		return
	}

	id := parts[1]

	if err := c.service.DeleteTask(id); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("✓ Tarea eliminada correctamente")
}
