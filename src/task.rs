use serde::{Deserialize, Serialize};
use std::time::{SystemTime, UNIX_EPOCH};
use chrono::{DateTime, Local, TimeZone, Utc};

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "lowercase")]
pub enum TaskStatus {
    Todo,
    InProgress,
    Done,
}

impl TaskStatus {
    pub fn as_str(&self) -> &str {
        match self {
            TaskStatus::Todo => "todo",
            TaskStatus::InProgress => "in-progress",
            TaskStatus::Done => "done",
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Task {
    pub id: u32,
    pub description: String,
    pub status: TaskStatus,
    pub created_at: u64,
    pub updated_at: u64,
}

impl Task {
    pub fn new(id: u32, description: String) -> Self {
        let now = current_timestamp();
        Task {
            id,
            description,
            status: TaskStatus::Todo,
            created_at: now,
            updated_at: now,
        }
    }

    pub fn update_description(&mut self, description: String) {
        self.description = description;
        self.updated_at = current_timestamp();
    }

    pub fn set_status(&mut self, status: TaskStatus) {
        self.status = status;
        self.updated_at = current_timestamp();
    }

    pub fn display(&self) {
        let created = format_timestamp(self.created_at);
        let updated = format_timestamp(self.updated_at);
        
        println!(
            "[{}] {} - {}",
            self.id,
            self.status.as_str(),
            self.description
        );
        println!("    Creada: {} | Actualizada: {}", created, updated);
    }
}

fn current_timestamp() -> u64 {
    SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .expect("El tiempo retrocedió")
        .as_secs()
}

fn format_timestamp(timestamp: u64) -> String {
    // Convertir timestamp Unix a DateTime UTC
    let datetime_utc = Utc.timestamp_opt(timestamp as i64, 0)
        .single()
        .expect("Timestamp inválido");
    
    // Convertir a hora local del sistema
    let datetime_local: DateTime<Local> = datetime_utc.with_timezone(&Local);
    
    // Formatear como YYYY-MM-DD HH:MM:SS
    datetime_local.format("%Y-%m-%d %H:%M:%S").to_string()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_new_task_creation() {
        let task = Task::new(1, "Test task".to_string());
        assert_eq!(task.id, 1);
        assert_eq!(task.description, "Test task");
        assert_eq!(task.status, TaskStatus::Todo);
        assert_eq!(task.created_at, task.updated_at);
    }

    #[test]
    fn test_update_description() {
        let mut task = Task::new(1, "Original".to_string());
        let original_updated_at = task.updated_at;
        
        std::thread::sleep(std::time::Duration::from_secs(1));
        task.update_description("Updated".to_string());
        
        assert_eq!(task.description, "Updated");
        assert!(task.updated_at > original_updated_at);
    }

    #[test]
    fn test_set_status() {
        let mut task = Task::new(1, "Test".to_string());
        
        task.set_status(TaskStatus::InProgress);
        assert_eq!(task.status, TaskStatus::InProgress);
        
        task.set_status(TaskStatus::Done);
        assert_eq!(task.status, TaskStatus::Done);
    }

    #[test]
    fn test_task_status_as_str() {
        assert_eq!(TaskStatus::Todo.as_str(), "todo");
        assert_eq!(TaskStatus::InProgress.as_str(), "in-progress");
        assert_eq!(TaskStatus::Done.as_str(), "done");
    }

    #[test]
    fn test_format_timestamp() {
        let timestamp = 1700000000u64; // 2023-11-14 22:13:20 UTC
        let formatted = format_timestamp(timestamp);
        // Solo verificamos que el formato sea correcto (YYYY-MM-DD HH:MM:SS)
        assert!(formatted.len() == 19);
        assert!(formatted.contains("-"));
        assert!(formatted.contains(":"));
    }
}
