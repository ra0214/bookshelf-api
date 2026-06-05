# 🚀 Instalación y Configuración - JWT Authentication

## ⚡ Instalación en 5 Pasos

### Paso 1: Verificar Librería JWT

```bash
go get github.com/golang-jwt/jwt/v5
go mod tidy
```

✅ **Ya está instalada**

---

### Paso 2: Configurar JWT_SECRET

**Opción A: Temporal (para desarrollo)**

PowerShell:
```powershell
$env:JWT_SECRET = "tu-clave-secreta-mínimo-24-caracteres"
go run main.go
```

Bash/Linux:
```bash
export JWT_SECRET="tu-clave-secreta-mínimo-24-caracteres"
go run main.go
```

---

**Opción B: Permanente (.env)**

Crear archivo `.env` en la raíz:
```env
JWT_SECRET=tu-clave-secreta-muy-larga-y-segura-minimo-32-caracteres-aleatorios
```

Usar librería godotenv (opcional):
```bash
go get github.com/joho/godotenv
```

En main.go:
```go
import "github.com/joho/godotenv"

func main() {
    godotenv.Load()
    middleware.InitAuthMiddleware()
    // ...
}
```

---

### Paso 3: Generar Clave Secreta Segura

**PowerShell:**
```powershell
$bytes = [byte[]]::new(32)
$rng = [System.Security.Cryptography.RNGCryptoServiceProvider]::new()
$rng.GetBytes($bytes)
[Convert]::ToBase64String($bytes)
```

**Bash/Linux:**
```bash
openssl rand -base64 32
```

**Python:**
```python
import secrets
print(secrets.token_urlsafe(32))
```

**Online:** https://www.random.org/strings/ (Base64)

---

### Paso 4: Compilar

```bash
cd /path/to/kanban
go build ./...
```

✅ **Sin errores**

---

### Paso 5: Ejecutar

```bash
go run main.go
```

Verás:
```
[Auth] Middleware de autenticación inicializado
Servidor iniciado en http://localhost:8080
```

---

## ✅ Verificación

### Test 1: Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"usuario@example.com","password":"password123"}'
```

Esperado (200 OK):
```json
{
  "message": "login exitoso",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {...}
}
```

---

### Test 2: Crear Tarea (con token)

```bash
TOKEN="eyJhbGciOiJIUzI1NiIs..."

curl -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","status":"por_hacer"}'
```

Esperado (201 Created)

---

### Test 3: Sin Token (debe fallar)

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Test"}'
```

Esperado (401 Unauthorized):
```json
{
  "error": "token no proporcionado",
  "code": "MISSING_TOKEN"
}
```

---

## 🔧 Troubleshooting

### Error: "JWT_SECRET not configured"

**Solución:**
```bash
# Verificar que está configurada
echo $env:JWT_SECRET          # PowerShell
echo $JWT_SECRET              # Bash
```

---

### Error: "Token invalid or expired"

**Causas posibles:**
1. JWT_SECRET cambió (regenera token)
2. Token expiró (login de nuevo)
3. Token corrupto (verifica formato "Bearer <token>")

**Solución:**
```bash
# Login nuevamente para obtener nuevo token
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"usuario@example.com","password":"password123"}'
```

---

### Error: "go: missing go.sum entry"

**Solución:**
```bash
go mod tidy
go get github.com/golang-jwt/jwt/v5
go mod tidy
```

---

## 📦 Estructura de Respuesta

### Login (200 OK)
```json
{
  "message": "login exitoso",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "Juan Pérez",
    "email": "usuario@example.com"
  }
}
```

### Crear Tarea (201 Created)
```json
{
  "message": "tarea creada correctamente",
  "id": 1,
  "task": {
    "id": 1,
    "user_id": 1,
    "title": "Test",
    "status": "por_hacer",
    "created_at": "2024-05-31T10:30:00Z"
  }
}
```

### Sin Autenticación (401 Unauthorized)
```json
{
  "error": "token no proporcionado",
  "code": "MISSING_TOKEN"
}
```

---

## 🌐 CORS (si accedes desde frontend)

En main.go:
```go
import "github.com/gin-contrib/cors"

func main() {
    r := gin.Default()
    
    // Configurar CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Authorization", "Content-Type"},
        ExposeHeaders:    []string{"Authorization"},
        AllowCredentials: true,
    }))
    
    // ... resto del código
}
```

Instalar:
```bash
go get github.com/gin-contrib/cors
```

---

## 🐳 Docker (Producción)

Crear `Dockerfile`:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o kanban-api .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/kanban-api .
EXPOSE 8080

ENV JWT_SECRET=${JWT_SECRET}
CMD ["./kanban-api"]
```

Build:
```bash
docker build -t kanban-api .
```

Run:
```bash
docker run \
  -e JWT_SECRET="tu-clave-secreta" \
  -p 8080:8080 \
  kanban-api
```

---

## 📝 Environment Variables Producción

`.env.production`:
```env
# JWT Configuration
JWT_SECRET=clave-super-segura-minimo-32-caracteres-aleatorios
JWT_EXPIRY_HOURS=24

# Database
DB_HOST=prod-db.example.com
DB_USER=user
DB_PASSWORD=secure_password
DB_NAME=kanban_prod

# Server
SERVER_PORT=8080
GIN_MODE=release

# Logging
LOG_LEVEL=info
```

---

## 🔒 Checklist de Seguridad

- [ ] JWT_SECRET configurado (mín 32 caracteres)
- [ ] JWT_SECRET en variable de entorno (NO en código)
- [ ] HTTPS habilitado en producción
- [ ] CORS configurado para dominios permitidos
- [ ] Validación de input en todos los endpoints
- [ ] Password nunca en logs
- [ ] Token nunca en logs (solo primeros 10 caracteres)
- [ ] Rate limiting en /login
- [ ] Database credentials en variables de entorno
- [ ] SSL/TLS certificados válidos

---

## 📚 Documentación Disponible

1. **JWT_AUTHENTICATION_GUIDE.md** - Guía completa y detallada
2. **JWT_TESTING_EXAMPLES.md** - Ejemplos de testing
3. **JWT_QUICK_REFERENCE.md** - Referencia rápida

---

## ✨ Listo para Comenzar

```bash
# 1. Configurar JWT_SECRET
$env:JWT_SECRET = "tu-clave-secreta"

# 2. Compilar
go build ./...

# 3. Ejecutar
go run main.go

# 4. Test en otra terminal
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

---

**¡Listo! 🚀 Autenticación JWT completamente funcional**
