package application

import (
	"bookshelf/src/books/domain"
)

type EditBook struct {
	db domain.IBooks
}

func NewEditBook(db domain.IBooks) *EditBook {
	return &EditBook{db: db}
}

// Execute procesa la lógica de negocio para actualizar los datos de un libro existente
func (eb *EditBook) Execute(id int32, title, author, description, status string, rating *int32, review *string) error {
	return eb.db.UpdateBook(id, title, author, description, status, rating, review)
}