# Ejemplos Prácticos de Testing - JWT Authentication

## 🧪 Testing del Sistema de Autenticación JWT

Este documento contiene ejemplos prácticos para probar los endpoints con JWT.

---

## 1. Variables de Entorno

Antes de ejecutar, configura el JWT_SECRET:

**En PowerShell:**
```powershell
$env:JWT_SECRET = "mi-clave-secreta-muy-larga-y-segura-24-caracteres"
go run main.go
```

**En bash (Linux/Mac):**
```bash
export JWT_SECRET="mi-clave-secreta-muy-larga-y-segura-24-caracteres"
go run main.go
```

**En .env (opcional, con godotenv):**
```env
JWT_SECRET=mi-clave-secreta-muy-larga-y-segura-24-caracteres
```

---

## 2. Colecciones de Postman/Insomnia

### Crear usuario de prueba

**Request:**
```
POST http://localhost:8080/user
Content-Type: application/json

{
  "name": "Juan Pérez",
  "email": "juan@example.com",
  "password": "password123"
}
```

**Response (201 Created):**
```json
{
  "message": "Usuario creado correctamente",
  "id": 1
}
```

---

### Login (Obtener Token)

**Request:**
```
POST http://localhost:8080/login
Content-Type: application/json

{
  "email": "juan@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "message": "login exitoso",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6Imp1YW5AZXhhbXBsZS5jb20iLCJuYW1lIjoiSnVhbiBQw6lyZXoiLCJleHAiOjE3MTc2ODY2NjEsImlhdCI6MTcxNzYwMDI2MSwibmJmIjoxNzE3NjAwMjYxLCJpc3MiOiJrYW5iYW4tYXBpIn0.5q9_Z4_X2K_L3M_N4O_P5Q_R6S_T7U_V8W_X9Y_Z0",
  "user": {
    "id": 1,
    "name": "Juan Pérez",
    "email": "juan@example.com",
    "created_at": "2024-05-31T10:30:00Z"
  }
}
```

**Guardar el token (en Postman/Insomnia):**
```
Variables:
- auth_token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

### Crear Tarea (Protegido)

**Request:**
```
POST http://localhost:8080/tasks
Authorization: Bearer {{auth_token}}
Content-Type: application/json

{
  "title": "Implementar landing page",
  "description": "Crear página de inicio responsiva",
  "status": "por_hacer",
  "priority": "alta",
  "due_date": "2024-06-15"
}
```

**Response (201 Created):**
```json
{
  "message": "tarea creada correctamente",
  "id": 1,
  "task": {
    "id": 1,
    "user_id": 1,
    "title": "Implementar landing page",
    "description": "Crear página de inicio responsiva",
    "status": "por_hacer",
    "priority": "alta",
    "due_date": "2024-06-15",
    "created_at": "2024-05-31T10:35:00Z",
    "updated_at": "2024-05-31T10:35:00Z"
  }
}
```

---

### Ver Mis Tareas (Protegido)

**Request:**
```
GET http://localhost:8080/tasks/mias
Authorization: Bearer {{auth_token}}
```

**Response (200 OK):**
```json
{
  "message": "tareas obtenidas correctamente",
  "count": 1,
  "tasks": [
    {
      "id": 1,
      "user_id": 1,
      "title": "Implementar landing page",
      "description": "Crear página de inicio responsiva",
      "status": "por_hacer",
      "priority": "alta",
      "due_date": "2024-06-15",
      "created_at": "2024-05-31T10:35:00Z",
      "updated_at": "2024-05-31T10:35:00Z"
    }
  ]
}
```

---

### Actualizar Tarea (Protegido)

**Request:**
```
PUT http://localhost:8080/tasks/1
Authorization: Bearer {{auth_token}}
Content-Type: application/json

{
  "title": "Implementar landing page v2",
  "description": "Crear página de inicio responsiva con animaciones",
  "status": "en_progreso",
  "priority": "alta",
  "due_date": "2024-06-20"
}
```

**Response (200 OK):**
```json
{
  "message": "tarea actualizada correctamente",
  "task": {
    "id": 1,
    "user_id": 1,
    "title": "Implementar landing page v2",
    "description": "Crear página de inicio responsiva con animaciones",
    "status": "en_progreso",
    "priority": "alta",
    "due_date": "2024-06-20",
    "created_at": "2024-05-31T10:35:00Z",
    "updated_at": "2024-05-31T10:40:00Z"
  }
}
```

---

### Obtener Tarea por ID (Protegido)

**Request:**
```
GET http://localhost:8080/tasks/1
Authorization: Bearer {{auth_token}}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "user_id": 1,
  "title": "Implementar landing page v2",
  "description": "Crear página de inicio responsiva con animaciones",
  "status": "en_progreso",
  "priority": "alta",
  "due_date": "2024-06-20",
  "created_at": "2024-05-31T10:35:00Z",
  "updated_at": "2024-05-31T10:40:00Z"
}
```

---

### Eliminar Tarea (Protegido)

**Request:**
```
DELETE http://localhost:8080/tasks/1
Authorization: Bearer {{auth_token}}
```

**Response (200 OK):**
```json
{
  "message": "tarea eliminada correctamente"
}
```

---

## 3. Casos de Error

### Error: Token No Proporcionado

**Request:**
```
POST http://localhost:8080/tasks
Content-Type: application/json

{
  "title": "Nueva tarea"
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "token no proporcionado",
  "code": "MISSING_TOKEN"
}
```

---

### Error: Formato de Token Inválido

**Request:**
```
POST http://localhost:8080/tasks
Authorization: InvalidFormat token123
Content-Type: application/json

{
  "title": "Nueva tarea"
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "formato de token inválido. Usar: Bearer <token>",
  "code": "INVALID_FORMAT"
}
```

---

### Error: Token Expirado

**Request:**
```
POST http://localhost:8080/tasks
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3MTc2MDAwMDAwfQ...
Content-Type: application/json

{
  "title": "Nueva tarea"
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "token inválido: token is expired",
  "code": "INVALID_TOKEN"
}
```

---

### Error: Credenciales Inválidas

**Request:**
```
POST http://localhost:8080/login
Content-Type: application/json

{
  "email": "usuario@example.com",
  "password": "contraseña_incorrecta"
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "credenciales inválidas"
}
```

---

## 4. Script de Testing Automatizado

### Python (requests + unittest)

```python
import requests
import unittest
import json
import os

BASE_URL = "http://localhost:8080"
JWT_SECRET = os.getenv("JWT_SECRET", "mi-clave-secreta-muy-larga")

class TestJWTAuthentication(unittest.TestCase):
    
    @classmethod
    def setUpClass(cls):
        """Crear usuario de prueba"""
        response = requests.post(f"{BASE_URL}/user", json={
            "name": "Test User",
            "email": "test@example.com",
            "password": "password123"
        })
        cls.user_id = response.json()["id"]
    
    def test_login_successful(self):
        """Test: Login exitoso"""
        response = requests.post(f"{BASE_URL}/login", json={
            "email": "test@example.com",
            "password": "password123"
        })
        
        self.assertEqual(response.status_code, 200)
        data = response.json()
        self.assertIn("token", data)
        self.assertIn("user", data)
        self.token = data["token"]
    
    def test_login_invalid_password(self):
        """Test: Password incorrecto"""
        response = requests.post(f"{BASE_URL}/login", json={
            "email": "test@example.com",
            "password": "wrong_password"
        })
        
        self.assertEqual(response.status_code, 401)
        data = response.json()
        self.assertIn("error", data)
    
    def test_create_task_with_token(self):
        """Test: Crear tarea con token válido"""
        # Primero obtener token
        login_response = requests.post(f"{BASE_URL}/login", json={
            "email": "test@example.com",
            "password": "password123"
        })
        token = login_response.json()["token"]
        
        # Crear tarea
        response = requests.post(
            f"{BASE_URL}/tasks",
            headers={"Authorization": f"Bearer {token}"},
            json={
                "title": "Test Task",
                "description": "Testing JWT",
                "status": "por_hacer",
                "priority": "media"
            }
        )
        
        self.assertEqual(response.status_code, 201)
        data = response.json()
        self.assertIn("task", data)
    
    def test_create_task_without_token(self):
        """Test: Crear tarea sin token"""
        response = requests.post(
            f"{BASE_URL}/tasks",
            json={
                "title": "Test Task",
                "description": "Testing JWT"
            }
        )
        
        self.assertEqual(response.status_code, 401)
        data = response.json()
        self.assertEqual(data["code"], "MISSING_TOKEN")
    
    def test_view_my_tasks(self):
        """Test: Ver mis tareas"""
        # Primero obtener token
        login_response = requests.post(f"{BASE_URL}/login", json={
            "email": "test@example.com",
            "password": "password123"
        })
        token = login_response.json()["token"]
        
        # Ver tareas
        response = requests.get(
            f"{BASE_URL}/tasks/mias",
            headers={"Authorization": f"Bearer {token}"}
        )
        
        self.assertEqual(response.status_code, 200)
        data = response.json()
        self.assertIn("tasks", data)
        self.assertIsInstance(data["tasks"], list)

if __name__ == "__main__":
    unittest.main()
```

**Ejecutar tests:**
```bash
python -m unittest test_jwt.py -v
```

---

### JavaScript (Fetch + Jest)

```javascript
const BASE_URL = "http://localhost:8080";

describe("JWT Authentication", () => {
  let token;
  let userId;

  beforeAll(async () => {
    // Crear usuario de prueba
    const createResponse = await fetch(`${BASE_URL}/user`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: "Test User",
        email: "test@example.com",
        password: "password123"
      })
    });
    const createData = await createResponse.json();
    userId = createData.id;
  });

  test("Login exitoso", async () => {
    const response = await fetch(`${BASE_URL}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email: "test@example.com",
        password: "password123"
      })
    });

    expect(response.status).toBe(200);
    const data = await response.json();
    expect(data).toHaveProperty("token");
    expect(data).toHaveProperty("user");
    token = data.token;
  });

  test("Login con credenciales inválidas", async () => {
    const response = await fetch(`${BASE_URL}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email: "test@example.com",
        password: "wrong_password"
      })
    });

    expect(response.status).toBe(401);
  });

  test("Crear tarea con token", async () => {
    const response = await fetch(`${BASE_URL}/tasks`, {
      method: "POST",
      headers: {
        "Authorization": `Bearer ${token}`,
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        title: "Test Task",
        description: "Testing JWT",
        status: "por_hacer",
        priority: "media"
      })
    });

    expect(response.status).toBe(201);
    const data = await response.json();
    expect(data).toHaveProperty("task");
  });

  test("Acceder a recurso protegido sin token", async () => {
    const response = await fetch(`${BASE_URL}/tasks/mias`);
    
    expect(response.status).toBe(401);
    const data = await response.json();
    expect(data.code).toBe("MISSING_TOKEN");
  });
});
```

**Ejecutar tests:**
```bash
npm test
```

---

## 5. Debugging

### Ver contenido del JWT (sin verificar firma)

**Online Tool:** https://jwt.io

**Pegar el token en la sección "Encoded" para ver los claims sin verificar**

---

### Logs en el servidor

El middleware imprime logs útiles:

```
[Auth] Middleware de autenticación inicializado
[Login] Login exitoso para usuario: juan@example.com (ID: 1)
[Auth] Usuario 1 autenticado correctamente
[CreateTask] Obteniendo tareas para usuario: 1
```

---

## 6. Configuración para Producción

### Variables de Entorno (.env.production)

```env
JWT_SECRET=super-clave-secreta-minimo-32-caracteres-aleatorios-muy-segura
JWT_EXPIRY_HOURS=24
```

### Dockerfile

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app

# Copiar código
COPY . .

# Instalar dependencias
RUN go mod download

# Compilar
RUN go build -o kanban-api .

# Exponer puerto
EXPOSE 8080

# Ejecutar con variables de entorno
ENV JWT_SECRET=${JWT_SECRET}
CMD ["./kanban-api"]
```

**Build:**
```bash
docker build -t kanban-api .
docker run -e JWT_SECRET="tu-clave-secreta" -p 8080:8080 kanban-api
```

---

## 📋 Checklist de Testing

- [ ] Login exitoso
- [ ] Login con credenciales inválidas
- [ ] Crear tarea con token válido
- [ ] Acceder a recurso protegido sin token
- [ ] Token expirado
- [ ] Formato de token inválido
- [ ] Ver mis tareas
- [ ] Actualizar tarea
- [ ] Eliminar tarea
- [ ] CORS funcionando correctamente

---

**¡Listo para producción!** 🚀
