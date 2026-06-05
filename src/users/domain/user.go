package domain

import "time"

type IUser interface {
	SaveUser(name string, email string, password string) error
	DeleteUser(id int32) error
	UpdateUser(id int32, name string, email string, password string) error
	GetAll() ([]User, error)
	GetUserByCredentials(email string) (*User, error)
	GetUserByID(id int32) (*User, error)
}

type User struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func NewUser(name string, email string, password string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,
	}
}

func (u *User) SetName(name string) {
	u.Name = name
}
