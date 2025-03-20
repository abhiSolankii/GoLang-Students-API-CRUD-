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

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare select statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to execute select statement: %w", err)
	}
	defer rows.Close()
	var students []types.Student

	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		students = append(students, student)
	}
	return students, nil
}

func (s *Sqlite) UpdateStudentById(id int64, student types.Student) (types.Student, error) {
	//first find student and return error if not found
	_, err := s.GetStudentById(id)
	if err != nil {
		return types.Student{}, err
	}
	stmt, err := s.Db.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
	if err != nil {
		return types.Student{}, fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(student.Name, student.Email, student.Age, id)

	if err != nil {
		return types.Student{}, fmt.Errorf("failed to execute update statement: %w", err)
	}

	return student, nil

}
func (s *Sqlite) DeleteStudentById(id int64) (int64, error) {
	//first find student and return error if not found
	_, err := s.GetStudentById(id)
	if err != nil {
		return 0, err
	}

	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return 0, fmt.Errorf("failed to execute delete statement: %w", err)
	}
	return id, nil
}
