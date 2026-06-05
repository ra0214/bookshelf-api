package middleware

import (
	"fmt"
	"bookshelf/src/config/security"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ContextKeys define las claves utilizadas para almacenar valores en el contexto de Gin
const (
	ContextKeyUserID = "userID"
	ContextKeyEmail  = "email"
	ContextKeyName   = "name"
)

var jwtService security.JWTService

// InitAuthMiddleware inicializa el servicio JWT
// Debe llamarse una sola vez al iniciar la aplicación
func InitAuthMiddleware() {
	jwtService = security.NewJWTService()
	log.Println("[Auth] Middleware de autenticación inicializado")
}

// AuthMiddleware valida el JWT token de autenticación y extrae el ID del usuario
// Retorna 401 Unauthorized si el token no es válido, está expirado o no se proporciona
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("[Auth] Token no proporcionado en endpoint: %s", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token no proporcionado",
				"code":  "MISSING_TOKEN",
			})
			c.Abort()
			return
		}

		// El token viene en formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("[Auth] Formato de token inválido: %s", authHeader[:min(len(authHeader), 20)])
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "formato de token inválido. Usar: Bearer <token>",
				"code":  "INVALID_FORMAT",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validar el token usando el servicio JWT
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			log.Printf("[Auth] Error al validar token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("token inválido: %v", err),
				"code":  "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		// Guardar información del usuario en el contexto
		// Esta información estará disponible en los handlers posteriores
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyEmail, claims.Email)
		c.Set(ContextKeyName, claims.Name)

		log.Printf("[Auth] Usuario %d autenticado correctamente", claims.UserID)

		c.Next()
	}
}

// min retorna el mínimo entre dos números
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
