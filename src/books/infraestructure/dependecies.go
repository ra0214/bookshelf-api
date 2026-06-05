package infraestructure

import (
	"bookshelf/src/books/application"
	"bookshelf/src/books/domain"
	"bookshelf/src/config/middleware"
	userDomain "bookshelf/src/users/domain"

	"github.com/gin-gonic/gin"
)

// Init inicializa la conexión a MySQL y arranca el Setup del enrutador de libros
func Init(r *gin.Engine, userRepo userDomain.IUser) {
	// ps debe implementar la interfaz domain.IBooks (tu repositorio MySQL de libros)
	ps := NewMySQL() 
	SetupRouter(ps, r)
}

// SetupRouter inyecta las dependencias siguiendo Clean Architecture y define los endpoints HTTP
func SetupRouter(repo domain.IBooks, r *gin.Engine) {
	// 1. Inicialización de Casos de Uso (Application) y Controladores (Infrastructure)
	createBook := application.NewCreateBook(repo)
	createBookController := NewCreateBookController(createBook, repo)

	viewBooks := application.NewViewBooks(repo)
	viewBooksController := NewViewBooksController(viewBooks)

	viewMyBooks := application.NewViewMyBooks(repo)
	viewMyBooksController := NewViewMyBooksController(viewMyBooks)

	editBookUseCase := application.NewEditBook(repo)
	editBookController := NewEditBookController(editBookUseCase, repo)

	deleteBookUseCase := application.NewDeleteBook(repo)
	deleteBookController := NewDeleteBookController(deleteBookUseCase)

	getBookByIDUseCase := application.NewGetBookByID(repo)
	getBookByIDController := NewGetBookByIDController(getBookByIDUseCase)

	// 2. Definición de Rutas del Kanban de Lecturas (BookShelf Club)
	// Todas las rutas que gestionan la biblioteca personal usan el AuthMiddleware de JWT
	r.POST("/books", middleware.AuthMiddleware(), createBookController.Execute)
	r.GET("/books", viewBooksController.Execute) // Vista global o catálogo general
	r.GET("/books/mios", middleware.AuthMiddleware(), viewMyBooksController.Execute) // Mis Libros
	r.PUT("/books/:id", middleware.AuthMiddleware(), editBookController.Execute)
	r.DELETE("/books/:id", middleware.AuthMiddleware(), deleteBookController.Execute) // Nota: Asegúrate de que el controlador se llame igual
	r.GET("/books/:id", middleware.AuthMiddleware(), getBookByIDController.Execute)
}