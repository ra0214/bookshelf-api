package infraestructure

import (
	"bookshelf/src/config/security"
	"bookshelf/src/users/application"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginUserController struct {
	useCase    *application.LoginUser
	jwtService security.JWTService
}

func NewLoginUserController(useCase *application.LoginUser) *LoginUserController {
	return &LoginUserController{
		useCase:    useCase,
		jwtService: security.NewJWTService(),
	}
}

type LoginRequestBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponseBody struct {
	Message string        `json:"message"`
	Token   string        `json:"token"`
	User    *UserResponse `json:"user"`
}

// UserResponse es la estructura de respuesta con datos públicos del usuario
type UserResponse struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at,omitempty"`
}

// Execute maneja el login del usuario
// Si las credenciales son correctas, genera un JWT y lo devuelve
// - Retorna 200 OK con token si el login es exitoso
// - Retorna 400 Bad Request si los datos son inválidos
// - Retorna 401 Unauthorized si las credenciales son incorrectas
// - Retorna 500 Internal Server Error si hay error al generar el token
func (lc *LoginUserController) Execute(c *gin.Context) {
	var body LoginRequestBody

	// Parsear y validar el JSON
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Printf("[Login] Error al parsear JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "datos de entrada inválidos",
			"details": err.Error(),
		})
		return
	}

	// Ejecutar caso de uso de login (verifica credenciales)
	user, err := lc.useCase.Execute(body.Email, body.Password)
	if err != nil {
		log.Printf("[Login] Fallo de autenticación para usuario %s: %v", body.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "credenciales inválidas",
		})
		return
	}

	// Generar JWT con los datos del usuario
	token, err := lc.jwtService.GenerateToken(user.ID, user.Email, user.Name)
	if err != nil {
		log.Printf("[Login] Error al generar token para usuario %d: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error al generar token de autenticación",
		})
		return
	}

	// Preparar respuesta con datos públicos del usuario
	userResponse := &UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	if !user.CreatedAt.IsZero() {
		userResponse.CreatedAt = user.CreatedAt.String()
	}

	log.Printf("[Login] Login exitoso para usuario: %s (ID: %d)", user.Email, user.ID)

	// Responder con el token y datos del usuario
	c.JSON(http.StatusOK, LoginResponseBody{
		Message: "login exitoso",
		Token:   token,
		User:    userResponse,
	})
}
