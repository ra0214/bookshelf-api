package infraestructure

import (
	"bookshelf/src/books/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ViewBooksController struct {
	useCase *application.ViewBooks
}

func NewViewBooksController(useCase *application.ViewBooks) *ViewBooksController {
	return &ViewBooksController{useCase: useCase}
}

// Execute obtiene y retorna la lista global de todos los libros
// - Se activa mediante la ruta GET /books (vista de administración o catálogo general)
// - Retorna 200 OK con un arreglo de libros bajo la llave "books"
// - Retorna 500 Internal Server Error si ocurre un fallo en el caso de uso
func (v *ViewBooksController) Execute(ctx *gin.Context) {
	books, err := v.useCase.Execute()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"books": books})
}