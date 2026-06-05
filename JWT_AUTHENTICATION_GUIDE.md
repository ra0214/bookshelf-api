# JWT Authentication en Kanban API - Guía Completa

## 📋 Tabla de Contenidos
1. [Instalación de Dependencias](#instalación-de-dependencias)
2. [Estructura de Archivos](#estructura-de-archivos)
3. [Componentes Implementados](#componentes-implementados)
4. [Cómo Funciona el Flujo](#cómo-funciona-el-flujo)
5. [Ejemplos de Uso](#ejemplos-de-uso)
6. [Seguridad y Mejores Prácticas](#seguridad-y-mejores-prácticas)

---

## 1. Instalación de Dependencias

### Librería JWT
Se requiere la librería oficial de JWT para Go:

```bash
go get github.com/golang-jwt/jwt/v5
go mod tidy
```

**Versión instalada:** `github.com/golang-jwt/jwt/v5`

---

## 2. Estructura de Archivos

### Nuevos archivos creados:

```
kanban/
├── src/
│   ├── config/
│   │   ├── security/
│   │   │   └── jwt_service.go          ✅ NUEVO - Servicio JWT centralizado
│   │   └── middleware/
│   │       └── auth_middleware.go      ✅ MEJORADO - Middleware de autenticación
│   ├── tasks/
│   │   └── infraestructure/
│   │       ├── context_helper.go       ✅ NUEVO - Helper para extraer datos del contexto
│   │       ├── CreateTask_Controller.go  ✅ MEJORADO - Usa context_helper
│   │       ├── ViewMyTask_Controller.go  ✅ MEJORADO - Usa context_helper
│   │       └── EditTask_Controller.go    ✅ MEJORADO - Usa context_helper
│   └── users/
│       └── infraestructure/
│           └── LoginUser_Controller.go   ✅ MEJORADO - Usa servicio JWT centralizado
└── main.go                              ✅ MEJORADO - Inicializa JWT middleware
```

---

## 3. Componentes Implementados

### 📦 JWT Service (`src/config/security/jwt_service.go`)

**Responsabilidad:** Encapsular toda la lógica de JWT en una interfaz centralizada.

**Interface:**
```go
type JWTService interface {
    GenerateToken(userID int32, email string, name string) (string, error)
    ValidateToken(tokenString string) (*JWTClaims, error)
    RefreshToken(tokenString string) (string, error)
}
```

**Características:**
- ✅ Lee la clave secreta desde variable de entorno `JWT_SECRET`
- ✅ Genera tokens válidos por 24 horas
- ✅ Valida y parsea tokens JWT
- ✅ Manejo de errores con mensajes descriptivos
- ✅ Algoritmo de firma: HS256

**Uso:**
```go
jwtService := security.NewJWTService()

// Generar token
token, err := jwtService.GenerateToken(userID, email, name)

// Validar token
claims, err := jwtService.ValidateToken(tokenString)
if err == nil {
    fmt.Println("Token válido, UserID:", claims.UserID)
}

// Renovar token
newToken, err := jwtService.RefreshToken(oldToken)
```

---

### 🔐 Authentication Middleware (`src/config/middleware/auth_middleware.go`)

**Responsabilidad:** Interceptar requests HTTP y validar JWT tokens.

**Flujo:**
1. ✅ Extrae el header `Authorization: Bearer <token>`
2. ✅ Valida el formato "Bearer <token>"
3. ✅ Valida el token usando JWTService
4. ✅ Inyecta userID, email, y name en el contexto de Gin
5. ✅ Retorna 401 Unauthorized si hay error

**Constantes del Contexto:**
```go
const (
    ContextKeyUserID = "userID"  // int32
    ContextKeyEmail  = "email"   // string
    ContextKeyName   = "name"    // string
)
```

**Inicialización (en main.go):**
```go
func main() {
    // Inicializar PRIMERO el middleware JWT
    middleware.InitAuthMiddleware()
    
    // ... resto del código
}
```

**Aplicación a una ruta:**
```go
// Ruta protegida con JWT
r.POST("/tasks", middleware.AuthMiddleware(), createTaskController.Execute)

// Ruta pública (sin middleware)
r.GET("/tasks", viewTasksController.Execute)
```

---

### 👤 Login Controller (`src/users/infraestructure/LoginUser_Controller.go`)

**Responsabilidad:** Autenticar usuarios y generar JWT tokens.

**Flujo:**
1. ✅ Recibe email y password en JSON
2. ✅ Valida credenciales usando `LoginUser` use case
3. ✅ Genera JWT con datos del usuario
4. ✅ Retorna token + datos públicos del usuario

**Request:**
```json
POST /login
Content-Type: application/json

{
    "email": "usuario@example.com",
    "password": "micontraseña123"
}
```

**Response (200 OK):**
```json
{
    "message": "login exitoso",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
        "id": 1,
        "name": "Juan Pérez",
        "email": "usuario@example.com",
        "created_at": "2024-05-31T10:30:00Z"
    }
}
```

**Response (401 Unauthorized):**
```json
{
    "error": "credenciales inválidas"
}
```

---

### 🎯 Context Helper (`src/tasks/infraestructure/context_helper.go`)

**Responsabilidad:** Extraer datos del contexto de forma segura y consistente.

**Funciones disponibles:**
```go
// Extrae userID del contexto
userID := GetUserIDFromContext(c)  // int32

// Extrae email del contexto
email := GetEmailFromContext(c)    // string

// Extrae nombre del contexto
name := GetNameFromContext(c)      // string
```

---

## 4. Cómo Funciona el Flujo

### 📍 Flujo de Autenticación Completo

```
┌─────────────────────────────────────────────────────────────────┐
│ 1. CLIENTE REALIZA LOGIN                                        │
├─────────────────────────────────────────────────────────────────┤
│ POST /login                                                     │
│ Content-Type: application/json                                  │
│                                                                 │
│ {                                                               │
│     "email": "usuario@example.com",                            │
│     "password": "micontraseña"                                 │
│ }                                                               │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 2. LOGIN CONTROLLER VALIDA CREDENCIALES                        │
├─────────────────────────────────────────────────────────────────┤
│ - Ejecuta LoginUser use case                                    │
│ - Busca usuario en BD por email                                │
│ - Verifica password con bcrypt/hash                            │
│ - Si es válido: continúa                                        │
│ - Si es inválido: retorna 401 Unauthorized                     │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 3. JWT SERVICE GENERA TOKEN                                    │
├─────────────────────────────────────────────────────────────────┤
│ Claims (payload):                                              │
│ {                                                               │
│     "user_id": 42,                                             │
│     "email": "usuario@example.com",                            │
│     "name": "Juan Pérez",                                      │
│     "exp": 1717686600,  // Expira en 24 horas                │
│     "iat": 1717600200,  // Emitido ahora                      │
│     "iss": "kanban-api"                                        │
│ }                                                               │
│                                                                 │
│ Firma: HMAC-SHA256(header.payload, JWT_SECRET)               │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 4. SERVIDOR RETORNA TOKEN AL CLIENTE                          │
├─────────────────────────────────────────────────────────────────┤
│ HTTP/1.1 200 OK                                                │
│                                                                 │
│ {                                                               │
│     "message": "login exitoso",                                │
│     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",    │
│     "user": {...}                                              │
│ }                                                               │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 5. CLIENTE ALMACENA TOKEN                                      │
├─────────────────────────────────────────────────────────────────┤
│ LocalStorage, SessionStorage, o variable de sesión             │
│ Token se usa en todas las peticiones futuras                   │
└─────────────────────────────────────────────────────────────────┘
```

### 📍 Flujo de Request Protegido

```
┌─────────────────────────────────────────────────────────────────┐
│ 6. CLIENTE REALIZA REQUEST PROTEGIDO                           │
├─────────────────────────────────────────────────────────────────┤
│ POST /tasks                                                     │
│ Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9... │
│ Content-Type: application/json                                  │
│                                                                 │
│ {                                                               │
│     "title": "Nueva tarea",                                    │
│     "description": "Descripción...",                           │
│     "status": "por_hacer",                                     │
│     "priority": "alta"                                         │
│ }                                                               │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 7. MIDDLEWARE VALIDA TOKEN                                     │
├─────────────────────────────────────────────────────────────────┤
│ AuthMiddleware():                                               │
│ 1. Extrae header "Authorization"                              │
│ 2. Verifica formato "Bearer <token>"                          │
│ 3. Llama a jwtService.ValidateToken()                         │
│ 4. Verifica firma HMAC-SHA256                                 │
│ 5. Verifica que no esté expirado                              │
│                                                                 │
│ Si es válido:                                                  │
│ - Extrae claims (user_id, email, name)                        │
│ - Inyecta en contexto: c.Set("userID", claims.UserID)        │
│ - Continúa con handler                                         │
│                                                                 │
│ Si es inválido:                                                │
│ - Retorna 401 Unauthorized                                     │
│ - Aborta request                                               │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 8. HANDLER PROCESA REQUEST AUTENTICADO                         │
├─────────────────────────────────────────────────────────────────┤
│ CreateTask_Controller:                                          │
│ 1. Extrae userID del contexto: GetUserIDFromContext(c)        │
│ 2. Parsea JSON del request                                     │
│ 3. Llama a CreateTask use case con userID                      │
│ 4. Retorna 201 Created con tarea creada                        │
└─────────────────────────────────────────────────────────────────┘
```

---

## 5. Ejemplos de Uso

### 🔓 Login (Obtener Token)

**cURL:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@example.com",
    "password": "micontraseña123"
  }'
```

**Python requests:**
```python
import requests
import json

response = requests.post('http://localhost:8080/login', json={
    'email': 'usuario@example.com',
    'password': 'micontraseña123'
})

token = response.json()['token']
print(f"Token obtenido: {token}")
```

**JavaScript fetch:**
```javascript
const response = await fetch('http://localhost:8080/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        email: 'usuario@example.com',
        password: 'micontraseña123'
    })
});

const data = await response.json();
const token = data.token;
localStorage.setItem('auth_token', token);
```

---

### 🔐 Usar Token en Request Protegido

**cURL:**
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Implementar login",
    "description": "Añadir autenticación JWT",
    "status": "en_progreso",
    "priority": "alta"
  }'
```

**Python requests:**
```python
import requests

token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
headers = {
    'Authorization': f'Bearer {token}',
    'Content-Type': 'application/json'
}

response = requests.post('http://localhost:8080/tasks', 
    headers=headers,
    json={
        'title': 'Nueva tarea',
        'description': 'Descripción',
        'status': 'por_hacer',
        'priority': 'media'
    }
)

print(response.json())
```

**JavaScript fetch:**
```javascript
const token = localStorage.getItem('auth_token');

const response = await fetch('http://localhost:8080/tasks', {
    method: 'POST',
    headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
    },
    body: JSON.stringify({
        title: 'Nueva tarea',
        description: 'Descripción',
        status: 'por_hacer',
        priority: 'media'
    })
});

const data = await response.json();
console.log(data);
```

---

### ❌ Error - Token Inválido

**Request sin token:**
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Tarea"}'
```

**Response (401 Unauthorized):**
```json
{
    "error": "token no proporcionado",
    "code": "MISSING_TOKEN"
}
```

---

### ❌ Error - Token Expirado

**Response (401 Unauthorized):**
```json
{
    "error": "token inválido: token is expired",
    "code": "INVALID_TOKEN"
}
```

---

## 6. Seguridad y Mejores Prácticas

### 🔒 Variable de Entorno para JWT_SECRET

**En desarrollo (.env):**
```env
JWT_SECRET=mi-clave-secreta-muy-larga-y-segura-24-caracteres
```

**En producción:**
```bash
export JWT_SECRET="clave-super-segura-mínimo-32-caracteres-aleatorios"
```

**En Go:**
```go
package config

import "os"

func GetJWTSecret() string {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        // En producción, NUNCA permitir clave por defecto
        panic("JWT_SECRET no configurado")
    }
    return secret
}
```

### ✅ Requisitos de Seguridad

- ✅ **JWT_SECRET:** Mínimo 32 caracteres en producción
- ✅ **HTTPS:** Siempre usar HTTPS en producción (nunca HTTP)
- ✅ **Token Storage (Frontend):**
  - ✅ **Seguro:** localStorage, sessionStorage (no accessibles desde XSS con HttpOnly)
  - ❌ **Inseguro:** Cookies sin HttpOnly flag
- ✅ **CORS:** Configurar adecuadamente para limitar orígenes
- ✅ **Refresh Tokens:** Implementar tokens de refresco para rotación de tokens

### 🛡️ Mejoras Futuras

1. **Refresh Token System:**
   ```go
   // Token corta duración (15 min)
   // Refresh Token larga duración (7 días)
   ```

2. **Token Blacklist / Revocation:**
   ```go
   // Mantener lista de tokens invalidados (logout)
   ```

3. **Rate Limiting:**
   ```go
   // Limitar intentos de login fallidos
   ```

4. **Two-Factor Authentication (2FA):**
   ```go
   // Requerir código adicional tras login
   ```

---

## 📌 Resumen de Cambios

| Archivo | Tipo | Descripción |
|---------|------|-------------|
| `src/config/security/jwt_service.go` | ✅ NUEVO | Servicio JWT centralizado |
| `src/config/middleware/auth_middleware.go` | ✅ MEJORADO | Middleware con servicio JWT |
| `src/tasks/infraestructure/context_helper.go` | ✅ NUEVO | Helper para extraer datos del contexto |
| `src/tasks/infraestructure/CreateTask_Controller.go` | ✅ MEJORADO | Usa context_helper |
| `src/tasks/infraestructure/ViewMyTask_Controller.go` | ✅ MEJORADO | Usa context_helper |
| `src/tasks/infraestructure/EditTask_Controller.go` | ✅ MEJORADO | Usa context_helper |
| `src/users/infraestructure/LoginUser_Controller.go` | ✅ MEJORADO | Usa servicio JWT centralizado |
| `main.go` | ✅ MEJORADO | Inicializa JWT middleware |

---

## ✅ Checklist de Implementación

- [x] Instalar librería JWT
- [x] Crear servicio JWT centralizado
- [x] Implementar middleware de autenticación
- [x] Actualizar Login controller
- [x] Crear helper de contexto
- [x] Actualizar controladores de tareas
- [x] Inicializar middleware en main.go
- [x] Validar compilación
- [x] Documentación completa

---

## 🚀 Próximos Pasos

1. Configurar JWT_SECRET en variable de entorno
2. Implementar Refresh Token system
3. Añadir Rate Limiting al endpoint de login
4. Configurar CORS adecuadamente
5. Implementar logout con token blacklist (opcional)
6. Añadir tests unitarios
7. Documentar con Swagger/OpenAPI

---

**Hecho con ❤️ - Arquitectura Limpia + Hexagonal**
