package application

import (
	"bookshelf/src/books/domain"
)

type CreateBook struct {
	db domain.IBooks
}

func NewCreateBook(db domain.IBooks) *CreateBook {
	return &CreateBook{db: db}
}

// Execute procesa la lógica de negocio para registrar un nuevo libro en la biblioteca
func (cb *CreateBook) Execute(userID int32, title, author, description, status string, rating *int32, review *string) (int32, error) {
	return cb.db.CreateBook(userID, title, author, description, status, rating, review)
}