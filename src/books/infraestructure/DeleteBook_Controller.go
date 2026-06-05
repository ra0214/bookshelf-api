package infraestructure

import (
	"bookshelf/src/books/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteBookController struct {
	useCase *application.DeleteBook
}

func NewDeleteBookController(useCase *application.DeleteBook) *DeleteBookController {
	return &DeleteBookController{useCase: useCase}
}

// Execute elimina un libro mediante su ID recibido por parámetro en la URL
// - Requiere autenticación JWT en las rutas (para asegurar que la petición es válida)
// - Retorna 200 OK si el libro fue eliminado correctamente
// - Retorna 400 Bad Request si el ID proporcionado no es un número válido
// - Retorna 500 Internal Server Error si ocurre un fallo en la base de datos
func (d *DeleteBookController) Execute(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de libro inválido"})
		return
	}

	err = d.useCase.Execute(int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Error al eliminar el libro",
			"detalles": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Libro eliminado correctamente"})
}