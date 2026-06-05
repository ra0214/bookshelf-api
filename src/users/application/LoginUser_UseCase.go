package application

import (
	"errors"
	"bookshelf/src/users/domain"

	"golang.org/x/crypto/bcrypt"
)

type LoginUser struct {
	db domain.IUser
}

func NewLoginUser(db domain.IUser) *LoginUser {
	return &LoginUser{
		db: db,
	}
}

func (lu *LoginUser) Execute(email string, password string) (*domain.User, error) {
	// Obtener usuario solo con el email
	user, err := lu.db.GetUserByCredentials(email)
	if err != nil {
		return nil, err
	}

	// Verificar la contraseña
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	return user, nil
}
