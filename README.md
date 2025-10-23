# Task Tracker CLI

Un gestor de tareas modular escrito en Rust que te permite administrar tus tareas desde la línea de comandos.

## Características

- ✅ Agregar, actualizar y eliminar tareas
- 📝 Marcar tareas como en progreso o completadas
- 📋 Listar tareas por estado (todo, in-progress, done)
- 💾 Almacenamiento persistente en JSON
- 🏗️ Arquitectura modular y mantenible
- 🧪 Suite completa de pruebas unitarias e integración

## Instalación

\`\`\`bash
cargo build --release
\`\`\`

El ejecutable estará disponible en `target/release/task-tracker`

## Uso

### Agregar una tarea
\`\`\`bash
cargo run -- add "Comprar leche"
\`\`\`

### Actualizar una tarea
\`\`\`bash
cargo run -- update 1 "Comprar leche y pan"
\`\`\`

### Eliminar una tarea
\`\`\`bash
cargo run -- delete 1
\`\`\`

### Marcar como en progreso
\`\`\`bash
cargo run -- mark-in-progress 1
\`\`\`

### Marcar como completada
\`\`\`bash
cargo run -- mark-done 1
\`\`\`

### Listar todas las tareas
\`\`\`bash
cargo run -- list
\`\`\`

### Listar tareas por estado
\`\`\`bash
cargo run -- list todo
cargo run -- list in-progress
cargo run -- list done
\`\`\`

## Arquitectura

El proyecto está organizado en módulos:

- **main.rs**: Punto de entrada de la aplicación
- **lib.rs**: Módulo raíz de la biblioteca
- **task.rs**: Define la estructura `Task` y el enum `TaskStatus`
- **storage.rs**: Maneja la persistencia de datos en JSON
- **cli.rs**: Procesa los comandos de la línea de comandos

## Pruebas

El proyecto incluye una suite completa de pruebas:

### Ejecutar todas las pruebas
\`\`\`bash
cargo test
\`\`\`

### Ejecutar pruebas con salida detallada
\`\`\`bash
cargo test -- --nocapture
\`\`\`

### Ejecutar pruebas específicas
\`\`\`bash
# Pruebas unitarias de task.rs
cargo test --lib task::tests

# Pruebas unitarias de storage.rs
cargo test --lib storage::tests

# Pruebas unitarias de cli.rs
cargo test --lib cli::tests

# Pruebas de integración
cargo test --test integration_tests
\`\`\`

### Cobertura de pruebas

Las pruebas cubren:
- ✅ Creación y manipulación de tareas
- ✅ Operaciones CRUD en el almacenamiento
- ✅ Persistencia de datos entre sesiones
- ✅ Validación de comandos CLI
- ✅ Manejo de errores
- ✅ Flujos completos de usuario

## Almacenamiento

Las tareas se guardan en un archivo `tasks.json` en el directorio actual con el siguiente formato:

\`\`\`json
{
  "tasks": [
    {
      "id": 1,
      "description": "Comprar leche",
      "status": "todo",
      "created_at": 1234567890,
      "updated_at": 1234567890
    }
  ],
  "next_id": 2
}
\`\`\`

## Dependencias

- **serde**: Serialización y deserialización de datos
- **serde_json**: Manejo de formato JSON
- **chrono**: Manejo de fechas y zonas horarias

## Licencia

MIT
