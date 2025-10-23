use crate::storage::Storage;
use crate::task::TaskStatus;
use std::io;

pub struct Cli {
    storage: Storage,
}

impl Cli {
    pub fn new() -> Self {
        Cli {
            storage: Storage::new(),
        }
    }

    pub fn run(&self, args: &[String]) -> io::Result<()> {
        if args.is_empty() {
            return Err(io::Error::new(
                io::ErrorKind::InvalidInput,
                "No se proporcionó ningún comando",
            ));
        }

        let command = &args[0];

        match command.as_str() {
            "add" => self.handle_add(&args[1..]),
            "update" => self.handle_update(&args[1..]),
            "delete" => self.handle_delete(&args[1..]),
            "mark-in-progress" => self.handle_mark_in_progress(&args[1..]),
            "mark-done" => self.handle_mark_done(&args[1..]),
            "list" => self.handle_list(&args[1..]),
            _ => Err(io::Error::new(
                io::ErrorKind::InvalidInput,
                format!("Comando desconocido: {}", command),
            )),
        }
    }

    fn handle_add(&self, args: &[String]) -> io::Result<()> {
        if args.is_empty() {
            return Err(io::Error::new(
                io::ErrorKind::InvalidInput,
                "Se requiere una descripción para agregar una tarea",
            ));
        }

        let description = args.join(" ");
        let task = self.storage.add_task(description)?;
        println!("Tarea agregada exitosamente (ID: {})", task.id);
        Ok(())
    }

    fn handle_update(&self, args: &[String]) -> io::Result<()> {
        if args.len() < 2 {
            return Err(io::Error::new(
                io::ErrorKind::InvalidInput,
                "Se requiere un ID y una descripción para actualizar una tarea",
            ));
        }

        let id = self.parse_id(&args[0])?;
        let description = args[1..].join(" ");
        self.storage.update_task(id, description)?;
        println!("Tarea {} actualizada exitosamente", id);
        Ok(())
    }

    fn handle_delete(&self, args: &[String]) -> io::Result<()> {
        if args.is_empty() {
            return Err(io::Error::new(
                io::ErrorKind::InvalidInput,
                "Se requiere un ID para eliminar una tarea",
            ));
        }

        let id = self.parse_id(&args[0])?;
        self.storage.delete_task(id)?;
        println!("Tarea {} eliminada exitosamente", id);
        Ok(())
    }

    fn handle_mark_in_progress(&self, args: &[String]) -> io::Result<()> {
        if args.is_empty() {
            return Err(io::Error::new(
                io::ErrorKind::InvalidInput,
                "Se requiere un ID para marcar una tarea como en progreso",
            ));
        }

        let id = self.parse_id(&args[0])?;
        self.storage.mark_in_progress(id)?;
        println!("Tarea {} marcada como en progreso", id);
        Ok(())
    }

    fn handle_mark_done(&self, args: &[String]) -> io::Result<()> {
        if args.is_empty() {
            return Err(io::Error::new(
                io::ErrorKind::InvalidInput,
                "Se requiere un ID para marcar una tarea como completada",
            ));
        }

        let id = self.parse_id(&args[0])?;
        self.storage.mark_done(id)?;
        println!("Tarea {} marcada como completada", id);
        Ok(())
    }

    fn handle_list(&self, args: &[String]) -> io::Result<()> {
        let tasks = if args.is_empty() {
            self.storage.list_all()?
        } else {
            let filter = &args[0];
            match filter.as_str() {
                "done" => self.storage.list_by_status(TaskStatus::Done)?,
                "todo" => self.storage.list_by_status(TaskStatus::Todo)?,
                "in-progress" => self.storage.list_by_status(TaskStatus::InProgress)?,
                _ => {
                    return Err(io::Error::new(
                        io::ErrorKind::InvalidInput,
                        format!("Filtro desconocido: {}", filter),
                    ))
                }
            }
        };

        if tasks.is_empty() {
            println!("No hay tareas para mostrar");
        } else {
            println!("\n📋 Lista de Tareas");
            println!("{}", "=".repeat(70));
            for task in tasks {
                task.display();
                println!("{}", "-".repeat(70));
            }
        }

        Ok(())
    }

    fn parse_id(&self, id_str: &str) -> io::Result<u32> {
        id_str.parse::<u32>().map_err(|_| {
            io::Error::new(
                io::ErrorKind::InvalidInput,
                format!("ID inválido: {}", id_str),
            )
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_parse_id_valid() {
        let cli = Cli::new();
        let result = cli.parse_id("42");
        assert!(result.is_ok());
        assert_eq!(result.unwrap(), 42);
    }

    #[test]
    fn test_parse_id_invalid() {
        let cli = Cli::new();
        let result = cli.parse_id("not_a_number");
        assert!(result.is_err());
    }

    #[test]
    fn test_parse_id_negative() {
        let cli = Cli::new();
        let result = cli.parse_id("-5");
        assert!(result.is_err());
    }

    #[test]
    fn test_run_no_command() {
        let cli = Cli::new();
        let result = cli.run(&[]);
        assert!(result.is_err());
    }

    #[test]
    fn test_run_unknown_command() {
        let cli = Cli::new();
        let args = vec!["unknown".to_string()];
        let result = cli.run(&args);
        assert!(result.is_err());
    }

    #[test]
    fn test_handle_add_no_description() {
        let cli = Cli::new();
        let result = cli.handle_add(&[]);
        assert!(result.is_err());
    }

    #[test]
    fn test_handle_update_insufficient_args() {
        let cli = Cli::new();
        let result = cli.handle_update(&["1".to_string()]);
        assert!(result.is_err());
    }

    #[test]
    fn test_handle_delete_no_id() {
        let cli = Cli::new();
        let result = cli.handle_delete(&[]);
        assert!(result.is_err());
    }

    #[test]
    fn test_handle_list_invalid_filter() {
        let cli = Cli::new();
        let args = vec!["invalid_filter".to_string()];
        let result = cli.handle_list(&args);
        assert!(result.is_err());
    }
}
