package application

import (
	"bookshelf/src/books/domain"
)

type ViewMyBooks struct {
	db domain.IBooks
}

func NewViewMyBooks(db domain.IBooks) *ViewMyBooks {
	return &ViewMyBooks{db: db}
}

// Execute procesa la lógica para obtener la lista de libros de un usuario específico
func (vmb *ViewMyBooks) Execute(userID int32) ([]domain.Book, error) {
	return vmb.db.GetBooksByUserID(userID)
}