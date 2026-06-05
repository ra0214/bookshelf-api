package application

import "bookshelf/src/books/domain"

type GetBookByID struct {
	db domain.IBooks
}

func NewGetBookByID(db domain.IBooks) *GetBookByID {
	return &GetBookByID{db: db}
}

// Execute procesa la lógica de negocio para obtener un único libro mediante su ID
func (gb *GetBookByID) Execute(id int32) (*domain.Book, error) {
	return gb.db.GetBookByID(id)
}