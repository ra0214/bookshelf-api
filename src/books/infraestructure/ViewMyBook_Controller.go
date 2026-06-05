package infraestructure

import (
	"bookshelf/src/books/application"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ViewMyBooksController struct {
	useCase *application.ViewMyBooks
}

func NewViewMyBooksController(useCase *application.ViewMyBooks) *ViewMyBooksController {
	return &ViewMyBooksController{useCase: useCase}
}

// Execute devuelve todos los libros del usuario autenticado
// - Requiere autenticación JWT (Bearer token)
// - El userID se extrae del contexto (inyectado por AuthMiddleware)
// - Retorna 200 OK con array de libros
// - Retorna 401 Unauthorized si no hay token válido
// - Retorna 500 Internal Server Error si hay error en la BD
func (v *ViewMyBooksController) Execute(ctx *gin.Context) {
	// Extraer el userID del contexto usando la función helper del JWT
	userID := GetUserIDFromContext(ctx)
	if userID == 0 {
		log.Println("[ViewMyBooks] userID no encontrado en contexto")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	log.Printf("[ViewMyBooks] Obteniendo libros para usuario: %d", userID)

	// Ejecutar caso de uso para obtener los libros exclusivos del usuario
	books, err := v.useCase.Execute(userID)
	if err != nil {
		log.Printf("[ViewMyBooks] Error al obtener libros para usuario %d: %v", userID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "error al obtener los libros",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "libros obtenidos correctamente",
		"count":   len(books),
		"books":   books,
	})
}