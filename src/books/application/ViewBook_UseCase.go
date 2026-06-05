package application

import (
	"bookshelf/src/books/domain"
)

type ViewBooks struct {
	db domain.IBooks
}

func NewViewBooks(db domain.IBooks) *ViewBooks {
	return &ViewBooks{db: db}
}

// Execute procesa la lógica de negocio para obtener la lista global de todos los libros
func (vb *ViewBooks) Execute() ([]domain.Book, error) {
	return vb.db.GetAllBooks()
}