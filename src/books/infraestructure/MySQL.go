package infraestructure

import (
	"database/sql"
	"fmt"
	"log"

	"bookshelf/src/books/domain"
	"bookshelf/src/config"
)

type MySQL struct {
	conn *config.Conn_MySQL
}

// Validación en tiempo de compilación para asegurar que MySQL implementa domain.IBooks
var _ domain.IBooks = (*MySQL)(nil)

// NewMySQL inicializa la estructura y la retorna tipada como la interfaz del dominio
func NewMySQL() domain.IBooks {
	conn := config.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQL{conn: conn}
}

func (mysql *MySQL) CreateBook(userID int32, title, author, description, status string, rating *int32, review *string) (int32, error) {
	query := `INSERT INTO books (user_id, title, author, description, status, rating, review) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := mysql.conn.ExecutePreparedQuery(query, userID, title, author, description, status, rating, review)
	if err != nil {
		return 0, fmt.Errorf("error al registrar el libro: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener ID del libro creado: %v", err)
	}
	return int32(id), nil
}

func (mysql *MySQL) RelaciónDePrueba() {} // Opcional, si manejas extensiones

func (mysql *MySQL) GetAllBooks() ([]domain.Book, error) {
	query := `SELECT id, user_id, title, author, description, status, rating, review, created_at, updated_at FROM books`
	rows, err := mysql.conn.FetchRows(query)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar consulta global de libros: %v", err)
	}
	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		var book domain.Book
		var rating sql.NullInt32
		var review sql.NullString
		// El escaneo captura de forma segura los valores nulos (NullInt32 y NullString) de la BD
		if err := rows.Scan(&book.ID, &book.UserID, &book.Title, &book.Author, &book.Description, &book.Status, &rating, &review, &book.CreatedAt, &book.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error al escanear fila de libros: %v", err)
		}
		
		if rating.Valid {
			book.Rating = &rating.Int32
		}
		if review.Valid {
			book.Review = &review.String
		}
		
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar filas de libros: %v", err)
	}
	return books, nil
}

func (mysql *MySQL) GetBookByID(id int32) (*domain.Book, error) {
	query := `SELECT id, user_id, title, author, description, status, rating, review, created_at, updated_at FROM books WHERE id = ?`
	row, err := mysql.conn.FetchRow(query, id)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar consulta del libro %d: %v", id, err)
	}

	var book domain.Book
	var rating sql.NullInt32
	var review sql.NullString
	if err := row.Scan(&book.ID, &book.UserID, &book.Title, &book.Author, &book.Description, &book.Status, &rating, &review, &book.CreatedAt, &book.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no se encontró ningún libro con ID: %d", id)
		}
		return nil, fmt.Errorf("error al escanear los datos del libro: %v", err)
	}
	
	if rating.Valid {
		book.Rating = &rating.Int32
	}
	if review.Valid {
		book.Review = &review.String
	}
	
	return &book, nil
}

func (mysql *MySQL) GetBooksByUserID(userID int32) ([]domain.Book, error) {
	query := `SELECT id, user_id, title, author, description, status, rating, review, created_at, updated_at FROM books WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := mysql.conn.FetchRows(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar consulta de la biblioteca del usuario: %v", err)
	}
	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		var book domain.Book
		var rating sql.NullInt32
		var review sql.NullString
		if err := rows.Scan(&book.ID, &book.UserID, &book.Title, &book.Author, &book.Description, &book.Status, &rating, &review, &book.CreatedAt, &book.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error al escanear libro del usuario: %v", err)
		}
		
		if rating.Valid {
			book.Rating = &rating.Int32
		}
		if review.Valid {
			book.Review = &review.String
		}
		
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar la biblioteca del usuario: %v", err)
	}
	return books, nil
}

func (mysql *MySQL) UpdateBook(id int32, title, author, description, status string, rating *int32, review *string) error {
	query := `UPDATE books SET title = ?, author = ?, description = ?, status = ?, rating = ?, review = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := mysql.conn.ExecutePreparedQuery(query, title, author, description, status, rating, review, id)
	if err != nil {
		return fmt.Errorf("error al actualizar los datos del libro: %v", err)
	}
	return nil
}

func (mysql *MySQL) DeleteBook(id int32) error {
	query := "DELETE FROM books WHERE id = ?"
	_, err := mysql.conn.ExecutePreparedQuery(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar el libro de la base de datos: %v", err)
	}
	return nil
}