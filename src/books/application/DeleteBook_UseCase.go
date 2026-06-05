package application

import (
	"bookshelf/src/books/domain"
)

type DeleteBook struct {
	db domain.IBooks
}

func NewDeleteBook(db domain.IBooks) *DeleteBook {
	return &DeleteBook{db: db}
}

func (db *DeleteBook) Execute(id int32) error {
	return db.db.DeleteBook(id)
}