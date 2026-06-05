package infraestructure

import (
	"log"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext extrae el userID del contexto HTTP de forma segura
// Soporta diferentes tipos numéricos comunes causados por la decodificación de JWT
func GetUserIDFromContext(c *gin.Context) int32 {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		log.Printf("[ContextHelper] userID no encontrado en contexto para: %s", c.Request.URL.Path)
		return 0
	}

	// Switch de tipos para prevenir errores de casteo de JWT (ej. float64, int, int32)
	switch v := userIDInterface.(type) {
	case int32:
		return v
	case int:
		return int32(v)
	case float64:
		return int32(v)
	default:
		log.Printf("[ContextHelper] Tipo de userID no soportado (%T) para: %s", userIDInterface, c.Request.URL.Path)
		return 0
	}
}

// GetEmailFromContext extrae el email del contexto HTTP
func GetEmailFromContext(c *gin.Context) string {
	emailInterface, exists := c.Get("email")
	if !exists {
		return ""
	}

	email, ok := emailInterface.(string)
	if !ok {
		return ""
	}

	return email
}

// GetNameFromContext extrae el nombre del contexto HTTP
func GetNameFromContext(c *gin.Context) string {
	nameInterface, exists := c.Get("name")
	if !exists {
		return ""
	}

	name, ok := nameInterface.(string)
	if !ok {
		return ""
	}

	return name
}