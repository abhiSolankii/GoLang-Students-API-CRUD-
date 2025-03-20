package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/abhiSolankii/students-api-go-lang/internal/config"
	"github.com/abhiSolankii/students-api-go-lang/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// Set connection pooling parameters (optional but recommended)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Context with timeout for executing the query
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		age INTEGER CHECK(age > 0)
	)`)

	if err != nil {
		db.Close() // Close the DB if initialization fails
		return nil, fmt.Errorf("failed to create students table: %w", err)
	}

	return &Sqlite{Db: db}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()
	// Execute the prepared statement
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, fmt.Errorf("failed to execute insert statement: %w", err)
	}
	// Get the last inserted ID
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last inserted ID: %w", err)
	}
	return lastId, nil

}
func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, fmt.Errorf("failed to prepare select statement: %w", err)
	}
	defer stmt.Close()
	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.ID, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with ID %d", id)
		}
		return types.Student{}, fmt.Errorf("failed to execute select statement: %w", err)
	}
	return student, nil

}
