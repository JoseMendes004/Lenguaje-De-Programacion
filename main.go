package main

import (
	"fmt"
	"os"

	"github.com/usuario/task-manager/cli"
	"github.com/usuario/task-manager/repository"
	"github.com/usuario/task-manager/service"
)

const dataFile = "tasks.json"

func main() {
	// Inicializar repositorio
	repo := repository.NewTaskRepository(dataFile)

	// Inicializar servicio
	taskService, err := service.NewTaskService(repo)
	if err != nil {
		fmt.Printf("Error al inicializar el servicio: %v\n", err)
		os.Exit(1)
	}

	// Inicializar y ejecutar CLI
	taskCLI := cli.NewCLI(taskService)
	taskCLI.Run()
}
