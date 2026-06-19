package infraestructure

import (
	"database/sql"
	"errors"
	"fmt"
	"bookshelf/src/config"
	"bookshelf/src/users/domain"
	"log"
)

type MySQL struct {
	conn *config.Conn_MySQL
}

var _ domain.IUser = (*MySQL)(nil)

func NewMySQL() domain.IUser {
	conn := config.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQL{conn: conn}
}

func (mysql *MySQL) SaveUser(name string, email string, password string) error {
	query := "INSERT INTO users (name, email, password) VALUES (?, ?, ?)"
	result, err := mysql.conn.ExecutePreparedQuery(query, name, email, password)
	if err != nil {
		return fmt.Errorf("error al ejecutar la consulta: %v", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Usuario creado correctamente: Name:%s Email:%s", name, email)
	} else {
		log.Println("[MySQL] - No se insertó ninguna fila")
	}
	return nil
}

func (mysql *MySQL) GetAll() ([]domain.User, error) {
	query := "SELECT id, name, email, password, created_at FROM users"
	rows, err := mysql.conn.FetchRows(query)
	if err != nil {
		return nil, fmt.Errorf("Error al ejecutar la consulta SELECT: %v", err)
	}
	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		user, err := scanUser(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("Error al escanear la fila: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterando sobre las filas: %v", err)
	}
	return users, nil
}

func (mysql *MySQL) UpdateUser(id int32, name string, email string, password string) error {
	query := "UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?"
	result, err := mysql.conn.ExecutePreparedQuery(query, name, email, password, id)
	if err != nil {
		return fmt.Errorf("Error al ejecutar la consulta: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Usuario actualizado correctamente: ID: %d Name:%s Email: %s", id, name, email)
	} else {
		log.Println("[MySQL] - No se actualizó ninguna fila")
	}
	return nil
}

func (mysql *MySQL) DeleteUser(id int32) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := mysql.conn.ExecutePreparedQuery(query, id)
	if err != nil {
		return fmt.Errorf("Error al ejecutar la consulta: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Usuario eliminado correctamente: ID: %d", id)
	} else {
		log.Println("[MySQL] - No se eliminó ninguna fila")
	}
	return nil
}

func scanUser(scan func(dest ...any) error) (domain.User, error) {
	var user domain.User
	var createdAt sql.NullTime

	if err := scan(&user.ID, &user.Name, &user.Email, &user.Password, &createdAt); err != nil {
		return domain.User{}, err
	}

	if createdAt.Valid {
		user.CreatedAt = createdAt.Time
	}

	return user, nil
}

func (mysql *MySQL) GetUserByCredentials(email string) (*domain.User, error) {
	query := "SELECT id, name, email, password FROM users WHERE email = ?"
	row, err := mysql.conn.FetchRow(query, email)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta: %v", err)
	}

	var user domain.User
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		log.Printf("[MySQL] Error al leer usuario por email %s: %v", email, err)
		return nil, fmt.Errorf("error al leer usuario: %v", err)
	}

	return &user, nil
}

func (mysql *MySQL) GetUserByID(id int32) (*domain.User, error) {
	query := "SELECT id, name, email, password, created_at FROM users WHERE id = ?"
	row, err := mysql.conn.FetchRow(query, id)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta: %v", err)
	}

	user, err := scanUser(row.Scan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		log.Printf("[MySQL] Error al leer usuario por id %d: %v", id, err)
		return nil, fmt.Errorf("error al leer usuario: %v", err)
	}

	return &user, nil
}
