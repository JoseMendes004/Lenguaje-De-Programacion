use crate::task::{Task, TaskStatus};
use serde::{Deserialize, Serialize};
use std::fs;
use std::io::{self, ErrorKind};
use std::path::Path;

const TASKS_FILE: &str = "tasks.json";

#[derive(Debug, Serialize, Deserialize)]
struct TasksData {
    tasks: Vec<Task>,
    next_id: u32,
}

impl Default for TasksData {
    fn default() -> Self {
        TasksData {
            tasks: Vec::new(),
            next_id: 1,
        }
    }
}

pub struct Storage {
    file_path: String,
}

impl Storage {
    pub fn new() -> Self {
        Storage {
            file_path: TASKS_FILE.to_string(),
        }
    }

    pub fn with_path(path: String) -> Self {
        Storage { file_path: path }
    }

    fn read_data(&self) -> io::Result<TasksData> {
        if !Path::new(&self.file_path).exists() {
            return Ok(TasksData::default());
        }

        let content = fs::read_to_string(&self.file_path)?;
        
        if content.trim().is_empty() {
            return Ok(TasksData::default());
        }

        serde_json::from_str(&content).map_err(|e| {
            io::Error::new(ErrorKind::InvalidData, format!("Error al parsear JSON: {}", e))
        })
    }

    fn write_data(&self, data: &TasksData) -> io::Result<()> {
        let json = serde_json::to_string_pretty(data).map_err(|e| {
            io::Error::new(ErrorKind::InvalidData, format!("Error al serializar JSON: {}", e))
        })?;
        
        fs::write(&self.file_path, json)
    }

    pub fn add_task(&self, description: String) -> io::Result<Task> {
        let mut data = self.read_data()?;
        let task = Task::new(data.next_id, description);
        data.tasks.push(task.clone());
        data.next_id += 1;
        self.write_data(&data)?;
        Ok(task)
    }

    pub fn update_task(&self, id: u32, description: String) -> io::Result<()> {
        let mut data = self.read_data()?;
        
        let task = data.tasks.iter_mut().find(|t| t.id == id).ok_or_else(|| {
            io::Error::new(ErrorKind::NotFound, format!("Tarea con ID {} no encontrada", id))
        })?;

        task.update_description(description);
        self.write_data(&data)
    }

    pub fn delete_task(&self, id: u32) -> io::Result<()> {
        let mut data = self.read_data()?;
        let initial_len = data.tasks.len();
        data.tasks.retain(|t| t.id != id);

        if data.tasks.len() == initial_len {
            return Err(io::Error::new(
                ErrorKind::NotFound,
                format!("Tarea con ID {} no encontrada", id),
            ));
        }

        self.write_data(&data)
    }

    pub fn mark_in_progress(&self, id: u32) -> io::Result<()> {
        self.update_status(id, TaskStatus::InProgress)
    }

    pub fn mark_done(&self, id: u32) -> io::Result<()> {
        self.update_status(id, TaskStatus::Done)
    }

    fn update_status(&self, id: u32, status: TaskStatus) -> io::Result<()> {
        let mut data = self.read_data()?;
        
        let task = data.tasks.iter_mut().find(|t| t.id == id).ok_or_else(|| {
            io::Error::new(ErrorKind::NotFound, format!("Tarea con ID {} no encontrada", id))
        })?;

        task.set_status(status);
        self.write_data(&data)
    }

    pub fn list_all(&self) -> io::Result<Vec<Task>> {
        let data = self.read_data()?;
        Ok(data.tasks)
    }

    pub fn list_by_status(&self, status: TaskStatus) -> io::Result<Vec<Task>> {
        let data = self.read_data()?;
        Ok(data.tasks.into_iter().filter(|t| t.status == status).collect())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;

    fn setup_test_storage() -> (Storage, String) {
        let test_file = format!("test_tasks_{}.json", std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)
            .unwrap()
            .as_nanos());
        let storage = Storage::with_path(test_file.clone());
        (storage, test_file)
    }

    fn cleanup_test_file(file_path: &str) {
        let _ = fs::remove_file(file_path);
    }

    #[test]
    fn test_add_task() {
        let (storage, test_file) = setup_test_storage();
        
        let task = storage.add_task("Test task".to_string()).unwrap();
        assert_eq!(task.id, 1);
        assert_eq!(task.description, "Test task");
        assert_eq!(task.status, TaskStatus::Todo);
        
        cleanup_test_file(&test_file);
    }

    #[test]
    fn test_add_multiple_tasks() {
        let (storage, test_file) = setup_test_storage();
        
        let task1 = storage.add_task("Task 1".to_string()).unwrap();
        let task2 = storage.add_task("Task 2".to_string()).unwrap();
        
        assert_eq!(task1.id, 1);
        assert_eq!(task2.id, 2);
        
        cleanup_test_file(&test_file);
    }

    #[test]
    fn test_update_task() {
        let (storage, test_file) = setup_test_storage();
        
        let task = storage.add_task("Original".to_string()).unwrap();
        storage.update_task(task.id, "Updated".to_string()).unwrap();
        
        let tasks = storage.list_all().unwrap();
        assert_eq!(tasks[0].description, "Updated");
        
        cleanup_test_file(&test_file);
    }

    #[test]
    fn test_update_nonexistent_task() {
        let (storage, test_file) = setup_test_storage();
        
        let result = storage.update_task(999, "Updated".to_string());
        assert!(result.is_err());
        
        cleanup_test_file(&test_file);
    }

    #[test]
    fn test_delete_task() {
        let (storage, test_file) = setup_test_storage();
        
        let task = storage.add_task("To delete".to_string()).unwrap();
        storage.delete_task(task.id).unwrap();
        
        let tasks = storage.list_all().unwrap();
        assert_eq!(tasks.len(), 0);
        
        cleanup_test_file(&test_file);
    }

    #[test]
    fn test_delete_nonexistent_task() {
        let (storage, test_file) = setup_test_storage();
        
        let result = storage.delete_task(999);
        assert!(result.is_err());
        
        cleanup_test_file(&test_file);
    }

    #[test]
    fn test_mark_in_progress() {
        let (storage, test_file) = setup_test_storage();
        
        let task = storage.add_task("Test".to_string()).unwrap();
        storage.mark_in_progress(task.id).unwrap();
        
        let tasks = storage.list_all().unwrap();
        assert_eq!(tasks[0].status, TaskStatus::InProgress);
        
        cleanup_test_file(&test_file);
    }

    #[test]
    fn test_mark_done() {
        let (storage, test_file) = setup_test_storage();
        
        let task = storage.add_task("Test".to_string()).unwrap();
        storage.mark_done(task.id).unwrap();
        
        let tasks = storage.list_all().unwrap();
        assert_eq!(tasks[0].status, TaskStatus::Done);
        
        cleanup_test_file(&test_file);
    }

    #[test]
    fn test_list_by_status() {
        let (storage, test_file) = setup_test_storage();
        
        let task1 = storage.add_task("Task 1".to_string()).unwrap();
        let task2 = storage.add_task("Task 2".to_string()).unwrap();
        let task3 = storage.add_task("Task 3".to_string()).unwrap();
        
        storage.mark_in_progress(task2.id).unwrap();
        storage.mark_done(task3.id).unwrap();
        
        let todo_tasks = storage.list_by_status(TaskStatus::Todo).unwrap();
        let in_progress_tasks = storage.list_by_status(TaskStatus::InProgress).unwrap();
        let done_tasks = storage.list_by_status(TaskStatus::Done).unwrap();
        
        assert_eq!(todo_tasks.len(), 1);
        assert_eq!(in_progress_tasks.len(), 1);
        assert_eq!(done_tasks.len(), 1);
        
        cleanup_test_file(&test_file);
    }

    #[test]
    fn test_persistence() {
        let test_file = format!("test_persistence_{}.json", std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)
            .unwrap()
            .as_nanos());
        
        {
            let storage = Storage::with_path(test_file.clone());
            storage.add_task("Persistent task".to_string()).unwrap();
        }
        
        {
            let storage = Storage::with_path(test_file.clone());
            let tasks = storage.list_all().unwrap();
            assert_eq!(tasks.len(), 1);
            assert_eq!(tasks[0].description, "Persistent task");
        }
        
        cleanup_test_file(&test_file);
    }
}
