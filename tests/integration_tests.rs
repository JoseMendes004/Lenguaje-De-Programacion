use rust_task_tracker::cli::Cli;
use std::fs;

fn setup_test() -> String {
    let test_file = format!(
        "test_integration_{}.json",
        std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)
            .unwrap()
            .as_nanos()
    );
    test_file
}

fn cleanup_test(file_path: &str) {
    let _ = fs::remove_file(file_path);
}

#[test]
fn test_full_task_lifecycle() {
    let _test_file = setup_test();
    let cli = Cli::new();

    // Agregar tarea
    let add_args = vec!["add".to_string(), "Integration test task".to_string()];
    let result = cli.run(&add_args);
    assert!(result.is_ok());

    // Listar tareas
    let list_args = vec!["list".to_string()];
    let result = cli.run(&list_args);
    assert!(result.is_ok());

    // Marcar como en progreso
    let progress_args = vec!["mark-in-progress".to_string(), "1".to_string()];
    let result = cli.run(&progress_args);
    assert!(result.is_ok());

    // Actualizar descripción
    let update_args = vec![
        "update".to_string(),
        "1".to_string(),
        "Updated integration test".to_string(),
    ];
    let result = cli.run(&update_args);
    assert!(result.is_ok());

    // Marcar como completada
    let done_args = vec!["mark-done".to_string(), "1".to_string()];
    let result = cli.run(&done_args);
    assert!(result.is_ok());

    // Listar tareas completadas
    let list_done_args = vec!["list".to_string(), "done".to_string()];
    let result = cli.run(&list_done_args);
    assert!(result.is_ok());

    // Eliminar tarea
    let delete_args = vec!["delete".to_string(), "1".to_string()];
    let result = cli.run(&delete_args);
    assert!(result.is_ok());

    cleanup_test("tasks.json");
}

#[test]
fn test_multiple_tasks_workflow() {
    let _test_file = setup_test();
    let cli = Cli::new();

    // Agregar múltiples tareas
    cli.run(&vec!["add".to_string(), "Task 1".to_string()])
        .unwrap();
    cli.run(&vec!["add".to_string(), "Task 2".to_string()])
        .unwrap();
    cli.run(&vec!["add".to_string(), "Task 3".to_string()])
        .unwrap();

    // Cambiar estados
    cli.run(&vec!["mark-in-progress".to_string(), "2".to_string()])
        .unwrap();
    cli.run(&vec!["mark-done".to_string(), "3".to_string()])
        .unwrap();

    // Listar por diferentes estados
    let result = cli.run(&vec!["list".to_string(), "todo".to_string()]);
    assert!(result.is_ok());

    let result = cli.run(&vec!["list".to_string(), "in-progress".to_string()]);
    assert!(result.is_ok());

    let result = cli.run(&vec!["list".to_string(), "done".to_string()]);
    assert!(result.is_ok());

    cleanup_test("tasks.json");
}

#[test]
fn test_error_handling() {
    let _test_file = setup_test();
    let cli = Cli::new();

    // Intentar actualizar tarea inexistente
    let result = cli.run(&vec![
        "update".to_string(),
        "999".to_string(),
        "Should fail".to_string(),
    ]);
    assert!(result.is_err());

    // Intentar eliminar tarea inexistente
    let result = cli.run(&vec!["delete".to_string(), "999".to_string()]);
    assert!(result.is_err());

    // Intentar marcar tarea inexistente
    let result = cli.run(&vec!["mark-done".to_string(), "999".to_string()]);
    assert!(result.is_err());

    cleanup_test("tasks.json");
}
