package security

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService define la interfaz para operaciones con JWT
type JWTService interface {
	GenerateToken(userID int32, email string, name string) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
	RefreshToken(tokenString string) (string, error)
}

// JWTClaims estructura personalizada para JWT con información del usuario
type JWTClaims struct {
	UserID int32  `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// jwtServiceImpl implementación del servicio JWT
type jwtServiceImpl struct {
	secretKey      []byte
	expiryDuration time.Duration
}

// NewJWTService crea una nueva instancia del servicio JWT
// Lee la clave secreta de la variable de entorno JWT_SECRET o usa una por defecto
func NewJWTService() JWTService {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		// En producción, SIEMPRE usar variable de entorno
		secretKey = "tu-clave-secreta-change-this-in-production-24-chars-min"
	}

	return &jwtServiceImpl{
		secretKey:      []byte(secretKey),
		expiryDuration: 24 * time.Hour, // Tokens válidos por 24 horas
	}
}

// GenerateToken crea un nuevo token JWT con los datos del usuario
// El token incluye el userID, email y nombre en los claims
func (s *jwtServiceImpl) GenerateToken(userID int32, email string, name string) (string, error) {
	if userID <= 0 {
		return "", errors.New("userID debe ser mayor a 0")
	}

	if email == "" {
		return "", errors.New("email no puede estar vacío")
	}

	now := time.Now()
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiryDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "kanban-api",
		},
	}

	// Crear token con algoritmo HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("error al firmar token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken valida un token JWT y extrae los claims
// Retorna los claims si el token es válido, o error si está expirado/inválido
func (s *jwtServiceImpl) ValidateToken(tokenString string) (*JWTClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token no puede estar vacío")
	}

	claims := &JWTClaims{}

	// Parsear y validar el token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verificar que el método de firma sea HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error al parsear token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}

// RefreshToken genera un nuevo token basado en uno existente válido
// Útil para mantener sesiones activas sin que expire el token anterior
func (s *jwtServiceImpl) RefreshToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("no se puede renovar token inválido: %w", err)
	}

	// Generar nuevo token con los mismos datos
	return s.GenerateToken(claims.UserID, claims.Email, claims.Name)
}
