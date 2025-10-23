use std::env;
use std::process;
use task_tracker::cli::Cli;

fn main() {
    let args: Vec<String> = env::args().collect();
    
    if args.len() < 2 {
        print_usage();
        process::exit(1);
    }

    let cli = Cli::new();
    
    if let Err(e) = cli.run(&args[1..]) {
        eprintln!("Error: {}", e);
        process::exit(1);
    }
}

fn print_usage() {
    println!("Task Tracker - Gestiona tus tareas desde la línea de comandos\n");
    println!("Uso:");
    println!("  task-tracker add <descripción>           - Agregar una nueva tarea");
    println!("  task-tracker update <id> <descripción>   - Actualizar una tarea");
    println!("  task-tracker delete <id>                 - Eliminar una tarea");
    println!("  task-tracker mark-in-progress <id>       - Marcar tarea como en progreso");
    println!("  task-tracker mark-done <id>              - Marcar tarea como completada");
    println!("  task-tracker list                        - Listar todas las tareas");
    println!("  task-tracker list done                   - Listar tareas completadas");
    println!("  task-tracker list todo                   - Listar tareas pendientes");
    println!("  task-tracker list in-progress            - Listar tareas en progreso");
}
