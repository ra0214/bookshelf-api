package infraestructure

import (
	"bookshelf/src/books/application"
	"bookshelf/src/books/domain"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EditBookController struct {
	useCase *application.EditBook
	repo    domain.IBooks
}

func NewEditBookController(useCase *application.EditBook, repo domain.IBooks) *EditBookController {
	return &EditBookController{useCase: useCase, repo: repo}
}

// EditBookRequestBody define los campos aceptados en el JSON para actualizar un libro
type EditBookRequestBody struct {
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Description string  `json:"description"`
	Status      string  `json:"status"` // 'por_leer' | 'leyendo' | 'leido'
	Rating      *int32  `json:"rating,omitempty"`
	Review      *string `json:"review,omitempty"`
}

// Execute actualiza los detalles de un libro existente
// - Requiere autenticación JWT (Bearer token)
// - Retorna 200 OK con los datos del libro actualizado
// - Retorna 400 Bad Request si el ID es inválido o los datos del JSON son incorrectos
// - Retorna 401 Unauthorized si no hay un token válido en el contexto
// - Retorna 500 Internal Server Error si ocurre un fallo en la base de datos
func (e *EditBookController) Execute(ctx *gin.Context) {
	// Extraer el bookID del parámetro de la ruta (/books/:id)
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[EditBook] ID de libro inválido: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de libro inválido",
		})
		return
	}

	// Parsear el cuerpo de la solicitud JSON proveniente de Flutter
	var body EditBookRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Printf("[EditBook] Error al parsear JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "datos de entrada inválidos",
			"details": err.Error(),
		})
		return
	}

	// Extraer userID del contexto (inyectado previamente por el AuthMiddleware del JWT)
	userID := GetUserIDFromContext(ctx)
	if userID == 0 {
		log.Println("[EditBook] userID no encontrado en contexto")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	// Ejecutar caso de uso de actualización pasándole las propiedades del libro
	err = e.useCase.Execute(int32(id), body.Title, body.Author, body.Description, body.Status, body.Rating, body.Review)
	if err != nil {
		log.Printf("[EditBook] Error al actualizar libro %d para usuario %d: %v", id, userID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "error al actualizar el libro",
			"details": err.Error(),
		})
		return
	}

	// Obtener el libro modificado desde el repositorio para retornar el estado actual en el JSON
	book, err := e.repo.GetBookByID(int32(id))
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "libro actualizado correctamente",
			"book":    book,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "libro actualizado correctamente",
		})
	}
}