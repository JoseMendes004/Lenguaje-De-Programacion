// server.js

const express = require('express');
const cors = require('cors'); 
const app = express();
const port = 3000; 

// Middleware
app.use(cors()); 
app.use(express.json()); 

// --- Base de Datos (Simulada) ---
let todos = [
    { id: 1, text: 'Aprender Vue 3', completed: false },
    { id: 2, text: 'Crear la API con Express', completed: true },
    { id: 3, text: 'Conectar frontend y backend', completed: true },
];

let users = []; 
let userId = 1;
// Usuario de prueba:
users.push({ id: userId++, name: 'Tester', email: 'test@example.com', password: 'password' }); 

// --- ENDPOINTS DE AUTENTICACIÓN ---

// POST /api/register: REGISTRO
app.post('/api/register', (req, res) => {
    const { name, email, password } = req.body;
    if (!name || !email || !password) {
        return res.status(400).json({ message: 'Todos los campos son obligatorios.' });
    }
    if (users.find(u => u.email === email)) {
        return res.status(409).json({ message: 'El correo electrónico ya está registrado.' });
    }

    const newUser = { id: userId++, name, email, password };
    users.push(newUser);
    
    const { password: _, ...userWithoutPass } = newUser;
    res.status(201).json({ message: 'Registro exitoso.', user: userWithoutPass });
});


// POST /api/login: INICIO DE SESIÓN
app.post('/api/login', (req, res) => {
    const { email, password } = req.body;
    const user = users.find(u => u.email === email);

    if (!user || user.password !== password) {
        return res.status(401).json({ message: 'Credenciales inválidas.' });
    }

    const { password: _, ...userWithoutPass } = user;
    res.status(200).json({ message: 'Login exitoso.', user: userWithoutPass });
});


// --- ENDPOINTS DE ToDos ---

// GET /api/todos: Obtener todas las tareas
app.get('/api/todos', (req, res) => {
    res.status(200).json(todos);
});

// Inicia el servidor
app.listen(port, () => {
    console.log(`API de ToDos y Auth escuchando en http://localhost:${port}`);
});