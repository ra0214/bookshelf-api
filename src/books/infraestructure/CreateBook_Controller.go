package infraestructure

import (
	"bookshelf/src/books/application"
	"bookshelf/src/books/domain"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateBookController struct {
	useCase *application.CreateBook
	repo    domain.IBooks
}

func NewCreateBookController(useCase *application.CreateBook, repo domain.IBooks) *CreateBookController {
	return &CreateBookController{useCase: useCase, repo: repo}
}

// CreateBookRequestBody define los campos esperados en el JSON de entrada desde Flutter
type CreateBookRequestBody struct {
	Title       string  `json:"title" binding:"required"`
	Author      string  `json:"author" binding:"required"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Rating      *int32  `json:"rating,omitempty"` // Soportará nulo si no se ha leído
	Review      *string `json:"review,omitempty"` // Soportará nulo si no se ha leído
}

// Execute crea un nuevo libro para el usuario autenticado
// - Requiere autenticación JWT (Bearer token)
// - El userID se extrae del contexto (inyectado por AuthMiddleware)
// - Retorna 201 Created con los datos del libro creado
// - Retorna 400 Bad Request si los datos son inválidos o faltan campos requeridos
// - Retorna 401 Unauthorized si no hay token válido
// - Retorna 500 Internal Server Error si hay error en la BD
func (c *CreateBookController) Execute(ctx *gin.Context) {
	var body CreateBookRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Printf("[CreateBook] Error al parsear JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "datos de entrada inválidos",
			"details": err.Error(),
		})
		return
	}

	// Extraer el userID del contexto usando tu función helper existente
	userID := GetUserIDFromContext(ctx)
	if userID == 0 {
		log.Println("[CreateBook] userID no encontrado en el contexto")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	// Ejecutar caso de uso de creación de libro
	id, err := c.useCase.Execute(userID, body.Title, body.Author, body.Description, body.Status, body.Rating, body.Review)
	if err != nil {
		log.Printf("[CreateBook] Error al crear libro para usuario %d: %v", userID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "error al registrar el libro",
			"details": err.Error(),
		})
		return
	}

	// Obtener el libro recién creado para devolverlo completo en la respuesta JSON
	book, err := c.repo.GetBookByID(id)
	if err == nil {
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "libro registrado correctamente",
			"id":      id,
			"book":    book,
		})
	} else {
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "libro registrado correctamente",
			"id":      id,
		})
	}
}