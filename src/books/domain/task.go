package domain

import "time"

type IBooks interface {
	CreateBook(userID int32, title, author, description, status string, rating *int32, review *string) (int32, error)
	GetAllBooks() ([]Book, error)
	GetBookByID(id int32) (*Book, error)
	GetBooksByUserID(userID int32) ([]Book, error) // Para que el usuario solo vea sus propios libros
	UpdateBook(id int32, title, author, description, status string, rating *int32, review *string) error
	DeleteBook(id int32) error
}

type Book struct {
	ID          int32     `json:"id"`
	UserID      int32     `json:"user_id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"` // 'por_leer' | 'leyendo' | 'leido'
	Rating      *int32    `json:"rating,omitempty"` // Puntero para soportar nulos (1 a 5 estrellas)
	Review      *string   `json:"review,omitempty"` // Puntero para soportar nulos
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// NewBook es el constructor para la entidad pura de un libro
func NewBook(userID int32, title, author, description, status string, rating *int32, review *string) *Book {
	// Estado por defecto si viene vacío
	if status == "" {
		status = "por_leer"
	}

	return &Book{
		UserID:      userID,
		Title:       title,
		Author:      author,
		Description: description,
		Status:      status,
		Rating:      rating,
		Review:      review,
	}
}