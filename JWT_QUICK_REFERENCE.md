# JWT Authentication - Quick Reference

## 🚀 Inicio Rápido

### 1. Configurar JWT_SECRET

```bash
# PowerShell
$env:JWT_SECRET = "tu-clave-secreta-min-24-chars"
go run main.go

# Bash
export JWT_SECRET="tu-clave-secreta-min-24-chars"
go run main.go
```

### 2. Flujo Básico

```
1. POST /login → Obtener token
   ↓
2. POST /tasks + Authorization: Bearer <token>
   ↓
3. Middleware valida → Inyecta userID en contexto
   ↓
4. Handler procesa con userID seguro
```

---

## 📁 Archivos Clave

| Archivo | Responsabilidad |
|---------|-----------------|
| `src/config/security/jwt_service.go` | Generar y validar JWT |
| `src/config/middleware/auth_middleware.go` | Validar token en requests |
| `src/tasks/infraestructure/context_helper.go` | Extraer datos del contexto |
| `src/users/infraestructure/LoginUser_Controller.go` | Generar token en login |

---

## 🔐 Rutas Protegidas vs Públicas

### Públicas (sin middleware)
```go
r.POST("/login", loginUserController.Execute)
r.POST("/user", createUserController.Execute)
r.GET("/tasks", viewTasksController.Execute)
```

### Protegidas (con middleware)
```go
r.POST("/tasks", middleware.AuthMiddleware(), createTaskController.Execute)
r.GET("/tasks/mias", middleware.AuthMiddleware(), viewMyTasksController.Execute)
r.PUT("/tasks/:id", middleware.AuthMiddleware(), editTaskController.Execute)
r.DELETE("/tasks/:id", middleware.AuthMiddleware(), deleteTaskController.Execute)
```

---

## 💻 Código en Handlers

### Extraer userID en un Controller

```go
func (c *CreateTaskController) Execute(ctx *gin.Context) {
    // FORMA FÁCIL - Usar helper
    userID := GetUserIDFromContext(ctx)
    
    if userID == 0 {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
        return
    }
    
    // Usar userID normalmente
    id, err := c.useCase.Execute(userID, ...)
    // ...
}
```

### Helper Disponibles

```go
userID := GetUserIDFromContext(c)    // int32
email := GetEmailFromContext(c)      // string
name := GetNameFromContext(c)        // string
```

---

## 🧪 Testing Rápido

### Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"pass123"}'
```

### Crear Tarea (con token)
```bash
TOKEN="eyJhbGciOiJIUzI1NiIs..."

curl -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Nueva tarea","status":"por_hacer"}'
```

### Sin Token (Error 401)
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Nueva tarea"}'

# Resultado: 401 Unauthorized
```

---

## 📋 Respuestas HTTP

| Endpoint | Método | Código | Descripción |
|----------|--------|--------|-------------|
| `/login` | POST | 200 | Login exitoso, retorna token |
| `/login` | POST | 401 | Credenciales inválidas |
| `/tasks` | POST | 201 | Tarea creada |
| `/tasks` | POST | 401 | Sin token o token inválido |
| `/tasks/:id` | PUT | 200 | Tarea actualizada |
| `/tasks/:id` | DELETE | 200 | Tarea eliminada |
| `/tasks/mias` | GET | 200 | Tareas del usuario |
| `/tasks/mias` | GET | 401 | Sin autenticación |

---

## 🔑 JWT Claims (Payload)

```json
{
  "user_id": 42,
  "email": "usuario@example.com",
  "name": "Juan Pérez",
  "exp": 1717686600,
  "iat": 1717600200,
  "nbf": 1717600200,
  "iss": "kanban-api"
}
```

---

## 🛡️ Seguridad

- ✅ JWT_SECRET mín 24 caracteres
- ✅ Usar HTTPS en producción
- ✅ Tokens válidos por 24 horas
- ✅ Firma HS256 (HMAC-SHA256)
- ✅ Token en Authorization header
- ✅ No enviar password en respuesta

---

## 📚 Documentación Completa

- `JWT_AUTHENTICATION_GUIDE.md` - Guía detallada
- `JWT_TESTING_EXAMPLES.md` - Ejemplos de testing

---

## ❓ FAQ

**P: ¿Dónde guardo el token en frontend?**
R: localStorage, sessionStorage, o variable de estado en React/Vue

**P: ¿Cómo hacer que el token no expire tan rápido?**
R: Modificar `s.expiryDuration = 24 * time.Hour` en jwt_service.go

**P: ¿Cómo implementar "Recuérdame" (remember me)?**
R: Usar Refresh Tokens con duraciones diferentes

**P: ¿Qué pasa si cambio JWT_SECRET?**
R: Todos los tokens anteriores se volverán inválidos

**P: ¿Puedo usar JWT_SECRET en el código?**
R: No, siempre usar variable de entorno en producción

---

**Última actualización:** 31 de mayo, 2024
