package main

import (
	"log"

	"bookshelf/src/config/middleware"
	booksInfra "bookshelf/src/books/infraestructure"
	usersInfra "bookshelf/src/users/infraestructure"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inicializar el middleware JWT
	middleware.InitAuthMiddleware()

	// 2. Inicializar base de datos / configuraciones de usuarios
	if err := usersInfra.InitUser(); err != nil {
		log.Fatalf("Error al inicializar usuarios: %v", err)
	}

	// 3. Crear el motor de enrutamiento de Gin
	r := gin.Default()

	// 4. Obtener el repositorio de usuarios (se necesita para el login/auth y pasarlo a libros)
	userRepo := usersInfra.NewMySQL()

	// 5. Configurar rutas del módulo de Usuarios (Login, Register)
	usersInfra.SetupRouter(userRepo, r)

	// 6. Inicializar y configurar rutas del módulo de Libros (BookShelf Club)
	// Le pasamos el router 'r' y el 'userRepo' tal como lo solicita tu arquitectura
	booksInfra.Init(r, userRepo)

	// 7. Iniciar el servidor HTTP
	log.Println("Servidor de BookShelf Club iniciado en http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}