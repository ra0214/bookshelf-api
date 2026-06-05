package application

import (
	"bookshelf/src/users/domain"
)

type EditUser struct {
	db domain.IUser
}

func NewEditUser(db domain.IUser) *EditUser {
	return &EditUser{db: db}
}

func (eu *EditUser) Execute(id int32, name string, email string, password string) error {
	return eu.db.UpdateUser(id, name, email, password)
}
