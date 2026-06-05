package infraestructure

import (
	"bookshelf/src/books/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetBookByIDController struct {
	useCase *application.GetBookByID
}

func NewGetBookByIDController(useCase *application.GetBookByID) *GetBookByIDController {
	return &GetBookByIDController{useCase: useCase}
}

// Execute obtiene los detalles de un libro específico mediante su ID
// - Se activa al consumir la ruta GET /books/:id
// - Retorna 200 OK con el objeto JSON del libro si se encuentra
// - Retorna 400 Bad Request si el parámetro ID no es un número válido
// - Retorna 500 Internal Server Error si ocurre un fallo en el caso de uso / base de datos
func (g *GetBookByIDController) Execute(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de libro inválido"})
		return
	}

	book, err := g.useCase.Execute(int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Error al obtener el libro",
			"detalles": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, book)
}